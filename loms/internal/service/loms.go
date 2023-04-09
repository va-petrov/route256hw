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

type OutboxMessage struct {
	MsgID   int64
	Key     string
	Message string
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
	AddOutbox(ctx context.Context, key string, message string) error
	GetOutbox(ctx context.Context) ([]OutboxMessage, error)
	DeleteOutbox(ctx context.Context, msgID int64) error
}

type NotificationsSender interface {
	SendNotification(ctx context.Context, msg OutboxMessage) error
}

type Service struct {
	LOMSRepo                  LOMSRepository
	TXMan                     TransactionManager
	NotificationsSender       NotificationsSender
	UnpayedOrdersJob          *jobs.Job
	StaleReservationsJob      *jobs.Job
	SendOrderNotificationsJob *jobs.Job
}

func New(lomsRepo LOMSRepository, txman TransactionManager, sender NotificationsSender) *Service {
	result := &Service{
		LOMSRepo:            lomsRepo,
		TXMan:               txman,
		NotificationsSender: sender,
	}
	result.UnpayedOrdersJob = jobs.NewJob("Unpayed orders", func(ctx context.Context) error {
		return result.UnpayedOrders(ctx)
	}, 30*time.Second)
	result.StaleReservationsJob = jobs.NewJob("Delete stale reservations", func(ctx context.Context) error {
		return result.StaleReservations(ctx)
	}, 60*time.Second)
	result.SendOrderNotificationsJob = jobs.NewJob("Send order notifications job", func(ctx context.Context) error {
		return result.SendOrderNotifications(ctx)
	}, 10*time.Second)
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
	err = m.SendOrderNotificationsJob.Run(ctx)
	if err != nil {
		result = errors.WithMessage(result, fmt.Sprintf("error starting job %v", m.SendOrderNotificationsJob.Name))
	}
	return result
}
