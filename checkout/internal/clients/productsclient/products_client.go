package productsclient

import (
	"context"
	"route256/checkout/internal/service"
	"route256/libs/clientwrapper"

	"github.com/pkg/errors"
)

type Client struct {
	url        string
	getProduct func(ctx context.Context, req ProductRequest) (*ProductInfoResponse, error)
	token      string
}

func New(url string, token string) *Client {
	return &Client{
		url:        url,
		getProduct: clientwrapper.New[ProductRequest, ProductInfoResponse](url + "/get_product"),
		token:      token,
	}
}

type ProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type ProductInfoResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, sku uint32) (service.Product, error) {
	request := ProductRequest{
		Token: c.token,
		SKU:   sku,
	}

	response, err := c.getProduct(ctx, request)
	if err != nil {
		return service.Product{}, errors.Wrap(err, "making loms.createOrder request")
	}

	return service.Product{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}
