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
	result := make([]Stock, len(DummyStocks.Stocks))
	for i, stock := range DummyStocks.Stocks {
		result[i] = Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		}
	}

	return result, nil
}
