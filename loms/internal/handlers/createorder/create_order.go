package createorder

import (
	"context"
	"log"
	"route256/loms/internal/service"
)

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Request struct {
	User  int64       `json:"user"`
	Items []OrderItem `json:"items"`
}

func (r Request) Validate() error {
	// TODO: implement
	return nil
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type Handler struct {
	Service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("createOrder: %+v", req)
	return Response{
		OrderID: 125,
	}, nil
}
