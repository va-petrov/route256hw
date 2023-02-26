package deletefromcart

import (
	"context"
	"log"
	"route256/checkout/internal/service"

	"github.com/pkg/errors"
)

type Handler struct {
	businessLogic *service.Service
}

func New(businessLogic *service.Service) *Handler {
	return &Handler{
		businessLogic: businessLogic,
	}
}

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
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
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("deleteFromCart: %+v", req)

	err := h.businessLogic.DeleteFromCart(ctx, req.User, req.SKU, req.Count)
	return Response{}, err
}
