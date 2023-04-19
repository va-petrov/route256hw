package productsclient

//go:generate sh -c "mkdir -p mocks && rm -rf mocks/client_minimock.go"
//go:generate minimock -i Client -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"route256/checkout/internal/config"
	"route256/checkout/internal/service/model"
	"route256/libs/cache"
	"route256/libs/limiter"
	log "route256/libs/logger"
	productServiceAPI "route256/product-service/pkg/product"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	cache         cache.Cache[uint32, model.Product]
}

func New(ctx context.Context, config config.ProductService) Client {
	conn, err := grpc.Dial(
		config.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal("failed to connect to product-service server", zap.Error(err))
	}

	cacheConfig := cache.Config{
		MaxSize: config.CacheConfig.MaxSize,
		TTL:     config.CacheConfig.TTL,
	}
	switch config.CacheConfig.Type {
	case "lru":
		cacheConfig.Type = cache.LRUCache
	case "lfu":
		cacheConfig.Type = cache.LFUCache
	default:
		cacheConfig.Type = cache.Simple
	}
	log.Debug("creating cache with config", zap.Any("cacheConfig", cacheConfig))
	productsCache, err := cache.NewCache[uint32, model.Product](ctx, cacheConfig)
	if err != nil {
		log.Error(ctx, "error creating cache", zap.Error(err))
	}

	return &client{
		productClient: productServiceAPI.NewProductServiceClient(conn),
		conn:          conn,
		token:         config.Token,
		rateLimiter:   limiter.New(ctx, time.Duration(int64(time.Second)/int64(config.RateLimit))),
		maxConcurrent: int(config.MaxConcurrent),
		cache:         productsCache,
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

var (
	CacheRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "products_cache",
		Name:      "requests_total",
	},
	)
	CacheHitsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "products_cache",
		Name:      "requests_hits_total",
	},
	)
	HistogramResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "products_cache",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
	)
	HistogramResponseHitTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "products_cache",
		Name:      "histogram_response_hit_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
	)
	HistogramResponseMissTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "products_cache",
		Name:      "histogram_response_miss_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
	)
)

func (c *client) GetProduct(ctx context.Context, sku uint32) (model.Product, error) {
	timeStart := time.Now()
	if c.cache != nil {
		CacheRequestsCounter.Inc()
		if result, ok := c.cache.Get(ctx, sku); ok {
			log.Debug("cache hit for SKU", zap.Uint32("SKU", sku))
			CacheHitsCounter.Inc()
			HistogramResponseHitTime.Observe(time.Since(timeStart).Seconds())
			return *result, nil
		}
		log.Debug("cache miss for SKU", zap.Uint32("SKU", sku))
	}
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

	result := model.Product{
		Name:  response.Name,
		Price: response.Price,
	}

	if c.cache != nil {
		log.Debug("saving response to cache for SKU", zap.Uint32("SKU", sku))
		if !c.cache.Set(ctx, sku, result) {
			log.Info("can't save products.GetProduct request result in cache for", zap.Uint32("SKU", sku))
		} else {
			log.Debug("response saved to cache for SKU", zap.Uint32("SKU", sku))
		}
		HistogramResponseMissTime.Observe(time.Since(timeStart).Seconds())
		HistogramResponseTime.Observe(time.Since(timeStart).Seconds())
	}
	return result, nil
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
