package loms_v1

import (
	"context"
	"route256/loms/internal/service"
	"route256/loms/pkg/loms_v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *loms_v1.CreateOrderRequest) (*loms_v1.CreateOrderResponse, error) {
	items := make([]service.Item, len(req.GetItems()))
	for i, item := range req.GetItems() {
		items[i].SKU = item.GetSku()
		items[i].Count = uint16(item.GetCount())
	}

	orderID, err := i.lomsService.CreateOrder(ctx, req.GetUser(), items)
	if err != nil {
		return nil, err
	}

	return &loms_v1.CreateOrderResponse{
		OrderID: orderID,
	}, nil
}
