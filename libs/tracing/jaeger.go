package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"route256/libs/logger"
)

func Init(serviceName string) {
	var cfg *config.Configuration

	cfg = &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	cfg, err := cfg.FromEnv()
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}

	_, err = cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}
