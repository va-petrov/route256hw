package loms_v1

import (
	"context"
	"route256/loms/pkg/loms_v1"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *loms_v1.CancelOrderRequest) (*loms_v1.CancelOrderResponse, error) {
	err := i.lomsService.CancelOrder(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	return &loms_v1.CancelOrderResponse{}, nil
}
