package productsclient

import (
	"context"
	"log"
	"route256/checkout/internal/config"
	"route256/checkout/internal/service"
	"route256/libs/limiter"
	productServiceAPI "route256/product-service/pkg/product"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	productClient productServiceAPI.ProductServiceClient
	conn          *grpc.ClientConn
	token         string
	rateLimiter   *limiter.Limiter
	maxConcurrent int
}

func New(ctx context.Context, config config.ProductService) *Client {
	conn, err := grpc.Dial(config.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to product-service server: %v", err)
	}

	return &Client{
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

func (c *Client) GetProduct(ctx context.Context, sku uint32) (service.Product, error) {
	select {
	case <-ctx.Done():
		return service.Product{}, errors.New("getProduct request cancelled")
	case t := <-c.rateLimiter.C:
		log.Printf("getProduct at time: %v", t)
	}
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

// GetProductsInfo параллельное выполнение запросов к productsService для заполнения информации о товарах в корзине
// Максимальное количество одновременных запросов задается через конфигурацию, параметр maxConcurrent для сервиса
// Если параметр равен 0, то все товары запрашиваются параллельно без ограничений
func (c *Client) GetProductsInfo(ctx context.Context, items []service.CartItem) error {
	var errs error
	taskSource := make(chan *service.CartItem)

	concurrency := c.maxConcurrent
	if c.maxConcurrent == 0 || len(items) <= c.maxConcurrent {
		concurrency = len(items)
	}
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(ctx context.Context, taskNum int) {
			defer wg.Done()
			log.Printf("task %v started", taskNum)
			for {
				select {
				case <-ctx.Done():
					return
				case item := <-taskSource:
					if item == nil {
						log.Printf("taskSource closed, task %v finishing", taskNum)
						return
					}
					log.Printf("task %v requesting info for sku %v", taskNum, item.SKU)
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

func (c *Client) Close() error {
	c.rateLimiter.Stop()
	return c.conn.Close()
}
