package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"route256/checkout/internal/api/checkout_v1"
	"route256/checkout/internal/clients/lomsclient"
	"route256/checkout/internal/clients/productsclient"
	"route256/checkout/internal/config"
	"route256/checkout/internal/repository/postgres"
	"route256/checkout/internal/service"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/interceptors"
	log "route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracing"
	"sync"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcPort    = flag.String("addr", ":8080", "port to listen")
	metricsPort = flag.String("metrics", ":7080", "port for metrics")
	develMode   = flag.Bool("devel", false, "development mode")
)

func main() {
	flag.Parse()

	log.Init(*develMode, zap.String("service", "checkout"))
	err := config.Init()
	if err != nil {
		log.Fatal("config init", zap.Error(err))
	}
	tracing.Init("checkout")

	lomsClient := lomsclient.New(config.ConfigData.Services.Loms)
	defer lomsClient.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	productsClient := productsclient.New(ctx, config.ConfigData.Services.ProductService)
	defer productsClient.Close()
	pool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect db", zap.Error(err))
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping db", zap.Error(err))
	}

	metricsServerDone := &sync.WaitGroup{}
	metricsServerDone.Add(1)
	metricsServer := &http.Server{
		Addr: *metricsPort,
	}

	go func(ctx context.Context) {
		defer metricsServerDone.Done()
		http.Handle("/metrics", metrics.New())

		log.Info("listening http for metrics", zap.String("addr", *metricsPort))
		if err := metricsServer.ListenAndServe(); err != nil {
			log.Error(ctx, "Error starting metrics handler", zap.Error(err))
		}
	}(ctx)

	cartRepo := postgres.NewCartRepo(pool)
	checkoutService := service.New(lomsClient, productsClient, cartRepo)

	checkout_v1.NewCheckoutV1(checkoutService)

	lis, err := net.Listen("tcp", *grpcPort)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				interceptors.LoggingInterceptor,
			),
		),
	)

	reflection.Register(s)
	desc.RegisterCheckoutServiceServer(s, checkout_v1.NewCheckoutV1(checkoutService))

	log.Info("server listening", zap.String("grpcAddr", *grpcPort))

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Error(ctx, "Error stopping metrics handler", zap.Error(err))
	}
	metricsServerDone.Wait()
}
