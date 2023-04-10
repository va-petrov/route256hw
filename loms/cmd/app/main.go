package main

import (
	"context"
	"flag"
	"net"
	"os"
	"route256/libs/interceptors"
	log "route256/libs/logger"
	"route256/libs/tracing"
	"route256/loms/internal/api/loms_v1"
	"route256/loms/internal/config"
	"route256/loms/internal/repository/postgres"
	"route256/loms/internal/repository/postgres/tranman"
	"route256/loms/internal/sender/kafka"
	"route256/loms/internal/service"
	desc "route256/loms/pkg/loms_v1"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcPort  = flag.String("addr", ":8081", "the port to listen")
	develMode = flag.Bool("devel", false, "development mode")
)

var brokers = []string{
	"kafka1:29091",
	"kafka2:29092",
	"kafka3:29093",
}

func main() {
	flag.Parse()

	log.Init(*develMode, zap.String("service", "loms"))

	err := config.Init()
	if err != nil {
		log.Fatal("config init", zap.Error(err))
	}

	tracing.Init("loms")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect db", zap.Error(err))
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping db", zap.Error(err))
	}

	txman := tranman.NewTransactionManager(pool)
	lomsRepo := postgres.NewLOMSRepo(txman)

	sender, err := kafka.NewSender(brokers, "orders")
	if err != nil {
		log.Fatal("error connecting to kafka", zap.Error(err))
	}
	lomsService := service.New(lomsRepo, txman, sender)
	err = lomsService.StartJobs(ctx)
	if err != nil {
		log.Fatal("error starting jobs", zap.Error(err))
	}

	loms_v1.NewLOMSV1(lomsService)

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
	desc.RegisterLOMSServiceServer(s, loms_v1.NewLOMSV1(lomsService))

	log.Info("server listening", zap.String("grpcPort", *grpcPort))

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
