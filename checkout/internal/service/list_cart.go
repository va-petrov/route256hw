package service

import (
	"context"

	"github.com/pkg/errors"
)

type CartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

func (m *Service) ListCart(ctx context.Context, user int64) ([]CartItem, error) {
	result := make([]CartItem, len(DummyCart.Items))
	for i, item := range DummyCart.Items {
		productInfo, err := m.ProductService.GetProduct(ctx, item.SKU)
		if err != nil {
			return nil, errors.WithMessage(err, "checking stocks")
		}
		result[i] = CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  productInfo.Name,
			Price: productInfo.Price,
		}
	}

	return result, nil
}
