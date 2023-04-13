package loms_v1

import (
	"context"
	"route256/loms/internal/service"
	"route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *loms_v1.CreateOrderRequest) (*loms_v1.CreateOrderResponse, error) {
	userID := req.GetUser()
	reqItems := req.GetItems()

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("userID", userID)
		span.SetTag("items", reqItems)
	}

	items := make([]service.Item, len(reqItems))
	for i, item := range reqItems {
		items[i].SKU = item.GetSku()
		items[i].Count = uint16(item.GetCount())
	}

	orderID, err := i.lomsService.CreateOrder(ctx, userID, items)
	if err != nil {
		return nil, err
	}

	return &loms_v1.CreateOrderResponse{
		OrderID: orderID,
	}, nil
}
