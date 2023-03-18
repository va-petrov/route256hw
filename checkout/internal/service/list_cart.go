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
		productInfo, err := m.ProductService.GetProduct(ctx, item.SKU)
		if err != nil {
			return Cart{}, errors.WithMessage(err, "getting product info")
		}
		result.Items[i] = CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  productInfo.Name,
			Price: productInfo.Price,
		}
		result.TotalPrice += uint32(item.Count) * productInfo.Price
	}

	return result, nil
}
