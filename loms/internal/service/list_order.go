package service

import (
	"context"
)

func (m *Service) ListOrder(ctx context.Context, orderID int64) (Order, error) {
	result := DummyOrder

	return result, nil
}
