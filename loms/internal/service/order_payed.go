package service

import (
	"context"
)

func (m *Service) OrderPayed(ctx context.Context, orderID int64) error {
	return m.TXMan.RunSerializable(ctx, func(ctxTX context.Context) error {
		order, err := m.LOMSRepo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		if order.Status != OrderStatusAwaitingPayment {
			return ErrIncorrectOrderState
		}

		reservations, err := m.LOMSRepo.GetReserves(ctx, orderID)
		if err != nil {
			return err
		}
		for _, reservation := range reservations {
			if err := m.LOMSRepo.ShipStock(ctx, reservation.SKU, reservation.WarehouseID, uint16(reservation.Count)); err != nil {
				return err
			}
		}
		if err := m.LOMSRepo.CancelReservationsForOrder(ctx, orderID); err != nil {
			return err
		}

		return m.LOMSRepo.SetStatusOrder(ctx, orderID, "payed")
	})
}
