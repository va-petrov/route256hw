package checkout_v1

import (
	"context"
	"route256/checkout/pkg/checkout_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) Purchase(ctx context.Context, req *checkout_v1.PurchaseRequest) (*checkout_v1.PurchaseResponse, error) {
	userID := req.GetUser()

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("userID", userID)
	}

	orderID, err := i.checkoutService.Purchase(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &checkout_v1.PurchaseResponse{OrderID: orderID}, nil
}
