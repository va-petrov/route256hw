package lomsclient

import (
	"context"
	"route256/checkout/internal/service"
	"route256/libs/clientwrapper"

	"github.com/pkg/errors"
)

type Client struct {
	url         string
	createOrder func(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
	stocks      func(ctx context.Context, req StocksRequest) (*StocksResponse, error)
}

func New(url string) *Client {
	return &Client{
		url:         url,
		createOrder: clientwrapper.New[CreateOrderRequest, CreateOrderResponse](url + "/createOrder"),
		stocks:      clientwrapper.New[StocksRequest, StocksResponse](url + "/stocks"),
	}
}

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type CreateOrderRequest struct {
	User  int64       `json:"user"`
	Items []OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderID int64 `json:"orderID"`
}

func (c *Client) CreateOrder(ctx context.Context, order service.Order) (int64, error) {
	request := CreateOrderRequest{
		User: order.User,
	}
	request.Items = make([]OrderItem, len(order.Items))
	for i, item := range order.Items {
		request.Items[i].SKU = item.SKU
		request.Items[i].Count = item.Count
	}

	response, err := c.createOrder(ctx, request)
	if err != nil {
		return -1, errors.Wrap(err, "making loms.createOrder request")
	}

	return response.OrderID, nil
}

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type StocksResponse struct {
	Stocks []StocksItem `json:"stocks"`
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]service.Stock, error) {
	request := StocksRequest{SKU: sku}

	response, err := c.stocks(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "making loms.stocks request")
	}

	stocks := make([]service.Stock, 0, len(response.Stocks))
	for _, stock := range response.Stocks {
		stocks = append(stocks, service.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return stocks, nil
}
