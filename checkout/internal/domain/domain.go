package domain

import "context"

type StocksChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type Model struct {
	stocksChecker StocksChecker
}

func New(stocksChecker StocksChecker) *Model {
	return &Model{
		stocksChecker: stocksChecker,
	}
}
