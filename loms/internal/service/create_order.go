package service

import (
	"context"
)

func (m *Service) CreateOrder(ctx context.Context, userID int64, items []Item) (int64, error) {
	order := Order{
		User:  userID,
		Items: make([]Item, len(items)),
	}
	for i, item := range items {
		order.Items[i] = Item{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}
	orderID, err := m.LOMSRepo.CreateOrder(ctx, order)
	if err != nil {
		return -1, err
	}

	for _, item := range items {
		stocks, err := m.LOMSRepo.GetStocks(ctx, item.SKU, true)
		if err != nil {
			return -1, err
		}
		counter := uint64(item.Count)
		for _, stock := range stocks {
			if stock.Count > counter {
				stock.Count = counter
			}
			if err := m.LOMSRepo.MakeReserve(ctx, orderID, item.SKU, stock.WarehouseID, stock.Count); err != nil {
				return -1, err
			}
			counter -= stock.Count
			if counter == 0 {
				break
			}
		}
		if counter > 0 {
			if err := m.LOMSRepo.SetStatusOrder(ctx, orderID, "failed"); err != nil {
				return -1, err
			}
		}
	}

	if err := m.LOMSRepo.SetStatusOrder(ctx, orderID, "awaiting payment"); err != nil {
		return -1, err
	}
	return orderID, nil
}
