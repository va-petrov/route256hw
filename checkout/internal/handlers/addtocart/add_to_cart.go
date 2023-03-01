package addtocart

import (
	"context"
	"log"
	"route256/checkout/internal/service"
	"route256/libs/validate"
)

type Handler struct {
	Service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (r Request) Validate() error {
	return validate.Combine(validate.User(r.User), validate.SKU(r.SKU))
}

type Response struct {
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("addToCart: %+v", req)

	err := h.Service.AddToCart(ctx, req.User, req.SKU, req.Count)
	return Response{}, err
}
