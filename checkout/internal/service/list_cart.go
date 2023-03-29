package service

import (
	"context"
	"route256/checkout/internal/service/model"

	"github.com/pkg/errors"
)

func (m *Service) ListCart(ctx context.Context, user int64) (model.Cart, error) {
	items, err := m.CartRepo.GetCart(ctx, user)
	if err != nil {
		return model.Cart{}, errors.WithMessage(err, "carts db")
	}
	result := model.Cart{
		Items: make([]model.CartItem, len(items)),
	}

	for i, item := range items {
		result.Items[i] = model.CartItem{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}
	err = m.ProductService.GetProductsInfo(ctx, result.Items)
	if err != nil {
		return model.Cart{}, err
	}
	for _, item := range result.Items {
		result.TotalPrice += uint32(item.Count) * item.Price
	}

	return result, nil
}
