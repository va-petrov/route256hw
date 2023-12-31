package logger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"net/http"

	"go.uber.org/zap"
)

var globalLogger *zap.Logger
var staticFields []zap.Field

func Init(devel bool, fields ...zap.Field) {
	globalLogger = New(devel)
	staticFields = fields
}

func New(devel bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if devel {
		logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel) /* TODO: поменять на подходящий для прода уровень логов */
		logger, err = cfg.Build()
	}
	if err != nil {
		panic(err)
	}

	return logger
}

func Middleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(
			"incoming http request",
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
		)

		next.ServeHTTP(w, r)
	})
}

func Debug(msg string, fields ...zap.Field) {
	fields = append(fields, staticFields...)
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	fields = append(fields, staticFields...)
	globalLogger.Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, staticFields...)
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if spancontext, ok := span.Context().(jaeger.SpanContext); ok {
			fields = append(
				fields,
				zap.String("trace", spancontext.TraceID().String()),
				zap.String("span", spancontext.SpanID().String()),
			)
		}
	}
	globalLogger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	fields = append(fields, staticFields...)
	globalLogger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, staticFields...)
	globalLogger.Fatal(msg, fields...)
}
