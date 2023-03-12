package service

import (
	"context"
)

func (m *Service) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	return m.CartRepo.DeleteFromCart(ctx, user, sku, count)
}
