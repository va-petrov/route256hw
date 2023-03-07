package main

import (
	"log"
	"net"
	"route256/checkout/internal/api/checkout_v1"
	"route256/checkout/internal/clients/lomsclient"
	"route256/checkout/internal/clients/productsclient"
	"route256/checkout/internal/config"
	"route256/checkout/internal/service"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/interceptors"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":8080"

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsClient := lomsclient.New(config.ConfigData.Services.Loms)
	defer lomsClient.Close()
	productsClient := productsclient.New(config.ConfigData.Services.ProductService.Url,
		config.ConfigData.Services.ProductService.Token)
	defer productsClient.Close()
	checkoutService := service.New(lomsClient, productsClient)

	checkout_v1.NewCheckoutV1(checkoutService)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
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

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
