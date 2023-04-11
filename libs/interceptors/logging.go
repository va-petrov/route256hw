package interceptors

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	log "route256/libs/logger"
	"route256/libs/metrics"
	"time"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Debug("incoming GRPC request", zap.String("method", info.FullMethod), zap.Any("request", req))
	metrics.RequestsCounter.WithLabelValues(info.FullMethod).Inc()

	timeStart := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			ext.Error.Set(span, true)
		}
		log.Error(ctx, "Error handling GRPC request", zap.String("method", info.FullMethod), zap.Error(err))
		metrics.ResponseCounter.WithLabelValues("error").Inc()

		elapsed := time.Since(timeStart)
		metrics.HistogramResponseTime.WithLabelValues("error").Observe(elapsed.Seconds())

		return nil, err
	}

	log.Debug("GRPC response", zap.String("method", info.FullMethod), zap.Any("response", res))
	metrics.ResponseCounter.WithLabelValues("success", info.FullMethod).Inc()

	elapsed := time.Since(timeStart)
	metrics.HistogramResponseTime.WithLabelValues("success", info.FullMethod).Observe(elapsed.Seconds())

	return res, nil
}
