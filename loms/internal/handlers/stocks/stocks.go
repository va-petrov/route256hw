package stocks

import (
	"context"
	"log"
	"route256/loms/internal/service"
)

type Request struct {
	SKU uint32 `json:"sku"`
}

func (r Request) Validate() error {
	// TODO: implement
	return nil
}

type Stock struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type Response struct {
	Stocks []Stock `json:"stocks"`
}

type Handler struct {
	Service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("stocks: %+v", req)

	stocks, err := h.Service.Stocks(ctx, req.SKU)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Stocks: make([]Stock, len(stocks)),
	}
	for i, stock := range stocks {
		response.Stocks[i] = Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		}
	}

	return response, nil
}
