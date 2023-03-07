package loms_v1

import (
	"context"
	"route256/loms/pkg/loms_v1"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *loms_v1.OrderPayedRequest) (*loms_v1.OrderPayedResponse, error) {
	err := i.lomsService.OrderPayed(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	return &loms_v1.OrderPayedResponse{}, nil
}
