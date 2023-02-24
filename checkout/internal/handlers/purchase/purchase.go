package purchase

import (
	"context"
	"log"
	"route256/checkout/internal/service"

	"github.com/pkg/errors"
)

type Handler struct {
	Services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{
		Services: services,
	}
}

type Request struct {
	User int64 `json:"user"`
}

var (
	ErrEmptyUser = errors.New("empty user")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	return nil
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("createOrder: %+v", req)

	var response Response

	orderID, err := h.Services.Purchase(ctx, req.User)
	if err != nil {
		return response, err
	}

	response.OrderID = orderID
	return response, nil
}
