package service

import (
	"context"
	"route256/checkout/internal/service/model"

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
	order := model.Order{
		User:  user,
		Items: items,
	}

	orderNo, err := m.LOMSService.CreateOrder(ctx, order)
	if err != nil {
		return -1, errors.WithMessage(err, "creating order")
	}

	if err := m.CartRepo.CleanCart(ctx, user); err != nil {
		return -1, errors.WithMessage(err, "cleaning cart")
	}

	return orderNo, nil
}
