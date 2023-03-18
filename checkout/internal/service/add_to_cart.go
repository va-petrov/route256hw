package service

import (
	"context"

	"github.com/pkg/errors"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func (m *Service) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	item, err := m.CartRepo.GetCartItem(ctx, user, sku)
	if err != nil {
		return errors.WithMessage(err, "carts db")
	}

	stocks, err := m.LOMSService.Stocks(ctx, sku)
	if err != nil {
		return errors.WithMessage(err, "checking stocks")
	}

	counter := int64(count)
	if item != nil {
		counter += int64(item.Count)
	}

	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			return m.CartRepo.AddToCart(ctx, user, sku, count)
		}
	}

	return ErrInsufficientStocks
}
