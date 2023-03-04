package checkout_v1

import (
	"context"
	"route256/pkg/checkout_v1"
)

func (i *Implementation) AddToCart(ctx context.Context, req *checkout_v1.AddToCartRequest) (*checkout_v1.AddToCartResponse, error) {
	err := i.checkoutService.AddToCart(ctx, req.GetUser(), req.GetSku(), uint16(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return &checkout_v1.AddToCartResponse{}, nil
}
