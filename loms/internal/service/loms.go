package service

import (
	"context"
	"fmt"
	"route256/libs/jobs"
	"time"

	"github.com/pkg/errors"
)

type Item struct {
	SKU   uint32
	Count uint16
}

const (
	OrderStatusFailed          = "failed"
	OrderStatusCancelled       = "cancelled"
	OrderStatusNew             = "new"
	OrderStatusAwaitingPayment = "awaiting payment"
	OrderStatusPayed           = "payed"
)

type Order struct {
	Status string
	User   int64
	Items  []Item
}

type Stock struct {
	SKU         uint32
	WarehouseID int64
	Count       uint64
}

type Stocks struct {
	Stocks []Stock
}

var (
	ErrIncorrectOrderState = errors.New("Incorrect order state for operation")
	ErrInsufficientStocks  = errors.New("insufficient stocks")
)

type TransactionManager interface {
	RunSerializable(ctx context.Context, fx func(ctxTX context.Context) error) error
	RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
	RunReadCommitted(ctx context.Context, fx func(ctxTX context.Context) error) error
	RunReadUncommitted(ctx context.Context, fx func(ctxTX context.Context) error) error
}

type LOMSRepository interface {
	GetStocks(ctx context.Context, sku uint32, checkReservations bool) ([]Stock, error)
	ShipStock(ctx context.Context, sku uint32, warehouseID int64, count uint16) error
	MakeReserve(ctx context.Context, orderID int64, sku uint32, warehouseID int64, count uint64) error
	GetReserves(ctx context.Context, orderID int64) ([]Stock, error)
	CancelReservationsForOrder(ctx context.Context, orderID int64) error
	CreateOrder(ctx context.Context, order Order) (int64, error)
	GetOrder(ctx context.Context, orderID int64) (*Order, error)
	SetStatusOrder(ctx context.Context, orderID int64, status string) error
	CancelUnpayedOrders(ctx context.Context) error
	DeleteStaleReservations(ctx context.Context) error
}

type Service struct {
	LOMSRepo             LOMSRepository
	TXMan                TransactionManager
	UnpayedOrdersJob     *jobs.Job
	StaleReservationsJob *jobs.Job
}

func New(lomsRepo LOMSRepository, txman TransactionManager) *Service {
	result := &Service{
		LOMSRepo: lomsRepo,
		TXMan:    txman,
	}
	result.UnpayedOrdersJob = jobs.NewJob("Unpayed orders", func(ctx context.Context) error {
		return result.UnpayedOrders(ctx)
	}, 30*time.Second)
	result.StaleReservationsJob = jobs.NewJob("Delete stale reservations", func(ctx context.Context) error {
		return result.StaleReservations(ctx)
	}, 60*time.Second)
	return result
}

func (m *Service) StartJobs(ctx context.Context) error {
	var result error
	err := m.UnpayedOrdersJob.Run(ctx)
	if err != nil {
		result = fmt.Errorf("error starting job %v", m.UnpayedOrdersJob.Name)
	}
	err = m.StaleReservationsJob.Run(ctx)
	if err != nil {
		result = errors.WithMessage(result, fmt.Sprintf("error starting job %v", m.StaleReservationsJob.Name))
	}
	return result
}
