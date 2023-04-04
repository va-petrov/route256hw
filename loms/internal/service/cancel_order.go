package service

import (
	"context"
	"fmt"
)

func (m *Service) CancelOrder(ctx context.Context, orderID int64) error {
	return m.TXMan.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		order, err := m.LOMSRepo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		if order.Status != OrderStatusAwaitingPayment {
			return ErrIncorrectOrderState
		}
		if err := m.LOMSRepo.CancelReservationsForOrder(ctx, orderID); err != nil {
			return err
		}

		if err := m.LOMSRepo.AddOutbox(ctx, fmt.Sprint(orderID), OrderStatusCancelled); err != nil {
			return err
		}

		return m.LOMSRepo.SetStatusOrder(ctx, orderID, "cancelled")
	})
}
