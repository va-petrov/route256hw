package service

import "context"

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
	WarehouseID int64
	Count       uint64
}

type Product struct {
	Name  string
	Price uint32
}

type LOMSClient interface {
	CreateOrder(ctx context.Context, order Order) (int64, error)
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type ProductClient interface {
	GetProduct(ctx context.Context, sku uint32) (Product, error)
}

type Service struct {
	LOMSService    LOMSClient
	ProductService ProductClient
}

func New(lomsClient LOMSClient, productClient ProductClient) *Service {
	return &Service{
		LOMSService:    lomsClient,
		ProductService: productClient,
	}
}

var DummyCart = Order{
	Items: []Item{
		{
			SKU:   1076963,
			Count: 10},
		{
			SKU:   1148162,
			Count: 5,
		},
		{
			SKU:   1625903,
			Count: 1,
		},
	},
}
