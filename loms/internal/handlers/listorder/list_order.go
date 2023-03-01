package listorder

import (
	"context"
	"log"
	"route256/libs/validate"
	"route256/loms/internal/service"
)

type Request struct {
	OrderID int64 `json:"orderID"`
}

func (r Request) Validate() error {
	return validate.OrderId(r.OrderID)
}

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Response struct {
	Status string      `json:"status"`
	User   int64       `json:"user"`
	Items  []OrderItem `json:"items"`
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
	log.Printf("listOrder: %+v", req)
	order, err := h.Service.ListOrder(ctx, req.OrderID)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Status: order.Status,
		User:   order.User,
		Items:  make([]OrderItem, len(order.Items)),
	}
	for i, item := range order.Items {
		response.Items[i] = OrderItem{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}

	return response, nil
}
