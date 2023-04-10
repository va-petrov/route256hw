package interceptors

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	log "route256/libs/logger"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Debug("incoming GRPC request", zap.String("method", info.FullMethod), zap.Any("request", req))

	res, err := handler(ctx, req)
	if err != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			ext.Error.Set(span, true)
		}
		log.Error(ctx, "Error handling GRPC request", zap.String("method", info.FullMethod), zap.Error(err))
		return nil, err
	}

	log.Debug("GRPC response", zap.String("method", info.FullMethod), zap.Any("response", res))

	return res, nil
}
