package cancelorder

import (
	"context"
	"log"
	"route256/loms/internal/service"
)

type Request struct {
	OrderID int64 `json:"orderID"`
}

func (r Request) Validate() error {
	// TODO: implement
	return nil
}

type Response struct {
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
	log.Printf("cancelOrder: %+v", req)
	err := h.Service.CancelOrder(ctx, req.OrderID)
	if err != nil {
		return Response{}, err
	}

	return Response{}, nil
}
