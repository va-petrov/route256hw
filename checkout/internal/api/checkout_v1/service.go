package checkout_v1

import (
	"route256/checkout/internal/service"
	"route256/pkg/checkout_v1"
)

type Implementation struct {
	checkout_v1.UnimplementedCheckoutServiceServer

	checkoutService *service.Service
}

func NewCheckoutV1(checkoutService *service.Service) *Implementation {
	return &Implementation{
		checkout_v1.UnimplementedCheckoutServiceServer{},
		checkoutService,
	}
}
