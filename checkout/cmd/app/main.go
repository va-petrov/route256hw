package main

import (
	"context"
	"flag"
	"net"
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

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcPort  = flag.String("addr", ":8080", "the port to listen")
	develMode = flag.Bool("devel", false, "development mode")
)

func main() {
	flag.Parse()

	log.Init(*develMode)
	err := config.Init()
	if err != nil {
		log.Fatal("config init", zap.Error(err))
	}

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
}
