package checkout_v1

import (
	"context"
	"route256/checkout/pkg/checkout_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *checkout_v1.DeleteFromCartRequest) (*checkout_v1.DeleteFromCartResponse, error) {
	userID := req.GetUser()
	sku := req.GetSku()
	count := uint16(req.GetCount())

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("userID", userID)
		span.SetTag("SKU", sku)
		span.SetTag("count", count)
	}

	err := i.checkoutService.DeleteFromCart(ctx, userID, sku, count)
	if err != nil {
		return nil, err
	}

	return &checkout_v1.DeleteFromCartResponse{}, nil
}
