package service

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

var (
	ErrEmptyCart = errors.New("Can't create order from empty cart")
)

func (m *Service) Purchase(ctx context.Context, user int64) (int64, error) {
	items, err := m.CartRepo.GetCart(ctx, user)
	if err != nil {
		return -1, errors.WithMessage(err, "getting cart from db")
	}
	if items == nil {
		return -1, ErrEmptyCart
	}
	order := Order{
		User:  user,
		Items: items,
	}

	orderNo, err := m.LOMSService.CreateOrder(ctx, order)
	if err != nil {
		return -1, errors.WithMessage(err, "creating order")
	}

	go func() {
		if err := m.CartRepo.CleanCart(context.Background(), user); err != nil {
			log.Printf("Error cleaning cart after creating order for user %v: %v", user, err)
		}
	}()

	return orderNo, nil
}
