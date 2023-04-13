package checkout_v1

import (
	"context"
	"route256/checkout/pkg/checkout_v1"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) ListCart(ctx context.Context, req *checkout_v1.ListCartRequest) (*checkout_v1.ListCartResponse, error) {
	userID := req.GetUser()

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("userID", userID)
	}

	cart, err := i.checkoutService.ListCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := checkout_v1.ListCartResponse{
		Items:      make([]*checkout_v1.CartItem, len(cart.Items)),
		TotalPrice: cart.TotalPrice,
	}
	for i, item := range cart.Items {
		response.Items[i] = &checkout_v1.CartItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		}
	}

	return &response, nil
}
