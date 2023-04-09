package productsclient

//go:generate sh -c "mkdir -p mocks && rm -rf mocks/client_minimock.go"
//go:generate minimock -i Client -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"route256/checkout/internal/config"
	"route256/checkout/internal/service/model"
	"route256/libs/limiter"
	log "route256/libs/logger"
	productServiceAPI "route256/product-service/pkg/product"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	GetProduct(ctx context.Context, sku uint32) (model.Product, error)
	GetProductsInfo(ctx context.Context, items []model.CartItem) error
	Close() error
}

type client struct {
	productClient productServiceAPI.ProductServiceClient
	conn          *grpc.ClientConn
	token         string
	rateLimiter   *limiter.Limiter
	maxConcurrent int
}

func New(ctx context.Context, config config.ProductService) Client {
	conn, err := grpc.Dial(config.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to product-service server", zap.Error(err))
	}

	return &client{
		productClient: productServiceAPI.NewProductServiceClient(conn),
		conn:          conn,
		token:         config.Token,
		rateLimiter:   limiter.New(ctx, time.Duration(int64(time.Second)/int64(config.RateLimit))),
		maxConcurrent: int(config.MaxConcurrent),
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

func (c *client) GetProduct(ctx context.Context, sku uint32) (model.Product, error) {
	select {
	case <-ctx.Done():
		return model.Product{}, errors.New("getProduct request cancelled")
	case t := <-c.rateLimiter.C:
		log.Debug("getProduct at time", zap.String("time", t.Format("2006-01-02 15:04:05.000000")))
	}
	request := productServiceAPI.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	}

	response, err := c.productClient.GetProduct(ctx, &request)
	if err != nil {
		return model.Product{}, errors.Wrap(err, "making loms.getProduct gRPC request")
	}

	return model.Product{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}

// GetProductsInfo параллельное выполнение запросов к productsService для заполнения информации о товарах в корзине
// Максимальное количество одновременных запросов задается через конфигурацию, параметр maxConcurrent для сервиса
// Если параметр равен 0, то все товары запрашиваются параллельно без ограничений
func (c *client) GetProductsInfo(ctx context.Context, items []model.CartItem) error {
	var errs error
	taskSource := make(chan *model.CartItem)

	concurrency := c.maxConcurrent
	if c.maxConcurrent == 0 || len(items) <= c.maxConcurrent {
		concurrency = len(items)
	}
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(ctx context.Context, taskNum int) {
			defer wg.Done()
			log.Debug("task started", zap.Int("taskNum", taskNum))
			for {
				select {
				case <-ctx.Done():
					return
				case item := <-taskSource:
					if item == nil {
						log.Debug("taskSource closed, task finishing", zap.Int("taskNum", taskNum))
						return
					}
					log.Debug("task requesting info for sku", zap.Int("taskNum", taskNum), zap.Uint32("SKU", item.SKU))
					product, err := c.GetProduct(ctx, item.SKU)
					if err != nil {
						if errs == nil {
							errs = err
						} else {
							errs = errors.WithMessage(errs, err.Error())
						}
					}
					item.Name = product.Name
					item.Price = product.Price
				}
			}
		}(ctx, i)
	}
	for i := range items {
		select {
		case <-ctx.Done():
			break
		case taskSource <- &items[i]:
		}
	}
	close(taskSource)
	wg.Wait()

	return errs
}

func (c *client) Close() error {
	c.rateLimiter.Stop()
	return c.conn.Close()
}
