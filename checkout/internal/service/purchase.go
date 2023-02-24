package service

import (
	"context"

	"github.com/pkg/errors"
)

func (m *Service) Purchase(ctx context.Context, user int64) (int64, error) {
	order := DummyCart
	order.User = user

	orderNo, err := m.LOMSService.CreateOrder(ctx, order)
	if err != nil {
		return -1, errors.WithMessage(err, "creating order")
	}
	return orderNo, nil
}
