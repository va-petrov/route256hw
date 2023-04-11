package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "requests_total",
	},
		[]string{"handler"},
	)
	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "responses_total",
	},
		[]string{"status", "handler"},
	)
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status", "handler"},
	)
)

func New() http.Handler {
	return promhttp.Handler()
}
