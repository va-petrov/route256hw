package checkout_v1

import (
	"context"
	"route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) Purchase(ctx context.Context, req *checkout_v1.PurchaseRequest) (*checkout_v1.PurchaseResponse, error) {
	orderID, err := i.checkoutService.Purchase(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	return &checkout_v1.PurchaseResponse{OrderID: orderID}, nil
}
