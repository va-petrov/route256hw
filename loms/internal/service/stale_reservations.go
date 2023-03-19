package service

import "context"

func (m *Service) StaleReservations(ctx context.Context) error {
	return m.LOMSRepo.DeleteStaleReservations(ctx)
}
