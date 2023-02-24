package lomsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"route256/checkout/internal/service"

	"github.com/pkg/errors"
)

type Client struct {
	url            string
	urlStocks      string
	urlCreateOrder string
}

func New(url string) *Client {
	return &Client{
		url:            url,
		urlStocks:      url + "/stocks",
		urlCreateOrder: url + "/createOrder",
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

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return -1, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlCreateOrder, bytes.NewBuffer(rawJSON))
	if err != nil {
		return -1, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return -1, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response CreateOrderResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return -1, errors.Wrap(err, "decoding json")
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

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlStocks, bytes.NewBuffer(rawJSON))
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response StocksResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json")
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
