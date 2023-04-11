package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	log "route256/libs/logger"
	"route256/libs/metrics"
	"route256/notifications/internal/kafka"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

var (
	metricsPort = flag.String("metrics", ":7082", "port for metrics")
	develMode   = flag.Bool("devel", false, "development mode")
)

var brokers = []string{
	"kafka1:29091",
	"kafka2:29092",
	"kafka3:29093",
}

func main() {
	flag.Parse()

	log.Init(*develMode, zap.String("service", "notifications"))

	ctx, cancel := context.WithCancel(context.Background())

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

	keepRunning := true
	log.Info("Starting notifications kafka consumer group...")

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	consumer := kafka.NewConsumerGroup()

	const groupName = "group-orders"

	client, err := sarama.NewConsumerGroup(brokers, groupName, config)
	if err != nil {
		log.Fatal("Error creating consumer group client", zap.Error(err))
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{"orders"}, &consumer); err != nil {
				log.Fatal("Error from consumer", zap.Error(err))
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready()
	log.Info("Notifications kafka consumer group ready...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Debug("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Debug("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Error(ctx, "Error stopping metrics handler", zap.Error(err))
	}
	metricsServerDone.Wait()

	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Fatal("Error closing client", zap.Error(err))
	}

}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Debug("Resuming consumption")
	} else {
		client.PauseAll()
		log.Debug("Pausing consumption")
	}

	*isPaused = !*isPaused
}
