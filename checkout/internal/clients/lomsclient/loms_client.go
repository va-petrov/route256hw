package lomsclient

//go:generate sh -c "mkdir -p mocks && rm -rf mocks/client_minimock.go"
//go:generate minimock -i Client -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"route256/checkout/internal/service/model"
	log "route256/libs/logger"
	lomsServiceAPI "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	CreateOrder(ctx context.Context, order model.Order) (int64, error)
	Stocks(ctx context.Context, sku uint32) ([]model.Stock, error)
	Close() error
}

type client struct {
	lomsClient lomsServiceAPI.LOMSServiceClient
	conn       *grpc.ClientConn
}

func New(url string) Client {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to loms server", zap.Error(err))
	}

	return &client{
		lomsClient: lomsServiceAPI.NewLOMSServiceClient(conn),
		conn:       conn,
	}
}

func (c *client) Close() error {
	return c.conn.Close()
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

func (c *client) CreateOrder(ctx context.Context, order model.Order) (int64, error) {
	request := lomsServiceAPI.CreateOrderRequest{
		User: order.User,
	}
	request.Items = make([]*lomsServiceAPI.OrderItem, len(order.Items))
	for i, item := range order.Items {
		request.Items[i] = &lomsServiceAPI.OrderItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		}
	}

	response, err := c.lomsClient.CreateOrder(ctx, &request)
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

func (c *client) Stocks(ctx context.Context, sku uint32) ([]model.Stock, error) {
	request := lomsServiceAPI.StocksRequest{Sku: sku}

	response, err := c.lomsClient.Stocks(ctx, &request)
	if err != nil {
		return nil, errors.Wrap(err, "making loms.stocks request")
	}

	stocks := make([]model.Stock, 0, len(response.Stocks))
	for _, stock := range response.Stocks {
		stocks = append(stocks, model.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return stocks, nil
}
