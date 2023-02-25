package listcart

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/service"
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

type CartItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Response struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"totalPrice"`
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("listCart: %+v", req)

	cart, err := h.Service.ListCart(ctx, req.User)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Items:      make([]CartItem, len(cart.Items)),
		TotalPrice: cart.TotalPrice,
	}
	for i, item := range cart.Items {
		response.Items[i] = CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		}
	}

	return response, nil
}
