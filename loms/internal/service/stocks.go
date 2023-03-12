package service

import (
	"context"
)

type CartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

func (m *Service) Stocks(ctx context.Context, sku uint32) ([]Stock, error) {
	stocks, err := m.LOMSRepo.GetStocks(ctx, sku, true)
	if err != nil {
		return nil, err
	}

	result := make([]Stock, len(stocks))
	for i, stock := range stocks {
		result[i] = Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		}
	}

	return result, nil
}
