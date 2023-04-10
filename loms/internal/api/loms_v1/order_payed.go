package loms_v1

import (
	"context"
	"route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *loms_v1.OrderPayedRequest) (*loms_v1.OrderPayedResponse, error) {
	orderID := req.GetOrderID()

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("orderID", orderID)
	}

	err := i.lomsService.OrderPayed(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &loms_v1.OrderPayedResponse{}, nil
}
