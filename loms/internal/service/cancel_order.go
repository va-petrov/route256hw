package service

import (
	"context"
)

func (m *Service) CancelOrder(ctx context.Context, orderID int64) error {
	order, err := m.LOMSRepo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status == "payed" {
		return ErrIncorrectOrderState
	}
	if err := m.LOMSRepo.CancelReservationsForOrder(ctx, orderID); err != nil {
		return err
	}
	return m.LOMSRepo.SetStatusOrder(ctx, orderID, "cancelled")
}
