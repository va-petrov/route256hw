package checkout_v1

import (
	"context"
	"route256/pkg/checkout_v1"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *checkout_v1.DeleteFromCartRequest) (*checkout_v1.DeleteFromCartResponse, error) {
	err := i.checkoutService.DeleteFromCart(ctx, req.GetUser(), req.GetSku(), uint16(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return &checkout_v1.DeleteFromCartResponse{}, nil
}
