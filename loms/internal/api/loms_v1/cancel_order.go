package loms_v1

import (
	"context"
	"route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *loms_v1.CancelOrderRequest) (*loms_v1.CancelOrderResponse, error) {
	orderID := req.GetOrderID()

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("orderID", orderID)
	}

	err := i.lomsService.CancelOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &loms_v1.CancelOrderResponse{}, nil
}
