package productsclient

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
	url           string
	urlGetProduct string
	token         string
}

func New(url string, token string) *Client {
	return &Client{
		url:           url,
		urlGetProduct: url + "/get_product",
		token:         token,
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

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return service.Product{}, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlGetProduct, bytes.NewBuffer(rawJSON))
	if err != nil {
		return service.Product{}, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return service.Product{}, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return service.Product{}, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response ProductInfoResponse

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return service.Product{}, errors.Wrap(err, "decoding json")
	}

	return service.Product{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}
