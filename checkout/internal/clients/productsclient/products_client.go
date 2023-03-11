package productsclient

import (
	"context"
	"log"
	"route256/checkout/internal/service"
	productServiceAPI "route256/product-service/pkg/product"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	productClient productServiceAPI.ProductServiceClient
	conn          *grpc.ClientConn
	token         string
}

func New(url string, token string) *Client {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to product-service server: %v", err)
	}

	return &Client{
		productClient: productServiceAPI.NewProductServiceClient(conn),
		conn:          conn,
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
	request := productServiceAPI.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	}

	response, err := c.productClient.GetProduct(ctx, &request)
	if err != nil {
		return service.Product{}, errors.Wrap(err, "making loms.getProduct gRPC request")
	}

	return service.Product{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
