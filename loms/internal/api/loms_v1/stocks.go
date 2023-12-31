package loms_v1

import (
	"context"
	"route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) Stocks(ctx context.Context, req *loms_v1.StocksRequest) (*loms_v1.StocksResponse, error) {
	sku := req.GetSku()

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("SKU", sku)
	}

	stocks, err := i.lomsService.Stocks(ctx, sku)
	if err != nil {
		return nil, err
	}

	response := loms_v1.StocksResponse{
		Stocks: make([]*loms_v1.StocksItem, len(stocks)),
	}
	for i, stock := range stocks {
		response.Stocks[i] = &loms_v1.StocksItem{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		}
	}

	return &response, nil
}
