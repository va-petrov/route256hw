package service

import "context"

func (m *Service) UnpayedOrders(ctx context.Context) error {
	return m.LOMSRepo.CancelUnpayedOrders(ctx)
}
