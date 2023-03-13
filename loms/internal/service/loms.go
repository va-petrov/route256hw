package service

import (
	"context"

	"github.com/pkg/errors"
)

type Item struct {
	SKU   uint32
	Count uint16
}

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
}

type Service struct {
	LOMSRepo LOMSRepository
	TXMan    TransactionManager
}

func New(lomsRepo LOMSRepository, txman TransactionManager) *Service {
	return &Service{
		LOMSRepo: lomsRepo,
		TXMan:    txman,
	}
}
