package service

import (
	"context"
)

func (m *Service) CreateOrder(ctx context.Context, userID int64, items []Item) (int64, error) {
	return 125, nil
}
