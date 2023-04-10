package loms_v1

import (
	"context"
	"route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) ListOrder(ctx context.Context, req *loms_v1.ListOrderRequest) (*loms_v1.ListOrderResponse, error) {
	orderID := req.GetOrderID()

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("orderID", orderID)
	}

	order, err := i.lomsService.ListOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	response := loms_v1.ListOrderResponse{
		Status: order.Status,
		User:   order.User,
		Items:  make([]*loms_v1.OrderItem, len(order.Items)),
	}
	for i, item := range order.Items {
		response.Items[i] = &loms_v1.OrderItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		}
	}
	return &response, nil
}
