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

type Cart struct {
	Items      []CartItem
	TotalPrice uint32
}

func (m *Service) ListCart(ctx context.Context, user int64) (Cart, error) {
	items, err := m.CartRepo.GetCart(ctx, user)
	if err != nil {
		return Cart{}, errors.WithMessage(err, "carts db")
	}
	result := Cart{
		Items: make([]CartItem, len(items)),
	}

	for i, item := range items {
		result.Items[i] = CartItem{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}
	err = m.ProductService.GetProductsInfo(ctx, result.Items)
	if err != nil {
		return Cart{}, err
	}
	for _, item := range result.Items {
		result.TotalPrice += uint32(item.Count) * item.Price
	}

	return result, nil
}
