package service

import (
	"context"
	"route256/checkout/internal/service/model"
)

type LOMSClient interface {
	CreateOrder(ctx context.Context, order model.Order) (int64, error)
	Stocks(ctx context.Context, sku uint32) ([]model.Stock, error)
}

type ProductClient interface {
	GetProduct(ctx context.Context, sku uint32) (model.Product, error)
	GetProductsInfo(ctx context.Context, items []model.CartItem) error
}

type CartRepository interface {
	GetCartItem(ctx context.Context, user int64, sku uint32) (*model.Item, error)
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error
	GetCart(ctx context.Context, user int64) ([]model.Item, error)
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
