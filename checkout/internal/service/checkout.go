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
	GetProductsInfo(ctx context.Context, items []CartItem) error
}

type CartRepository interface {
	GetCartItem(ctx context.Context, user int64, sku uint32) (*Item, error)
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error
	GetCart(ctx context.Context, user int64) ([]Item, error)
	CleanCart(ctx context.Context, user int64) error
}

type Service struct {
	LOMSService    LOMSClient
	ProductService ProductClient
	CartRepo       CartRepository
}

func New(lomsClient LOMSClient, productClient ProductClient, cartRepo CartRepository) *Service {
	return &Service{
		LOMSService:    lomsClient,
		ProductService: productClient,
		CartRepo:       cartRepo,
	}
}
