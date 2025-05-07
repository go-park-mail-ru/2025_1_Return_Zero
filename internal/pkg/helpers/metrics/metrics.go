package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	HTTPTotalNumberOfRequests *prometheus.CounterVec
	HTTPRequestDuration       *prometheus.HistogramVec
	// MicroserviceRequests        *prometheus.CounterVec
	// MicroserviceRequestDuration *prometheus.HistogramVec
	// MicroserviceErrors          *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer, namespace string) *Metrics {
	metrics := &Metrics{
		HTTPTotalNumberOfRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
				Namespace: namespace,
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:      "http_request_duration_seconds",
				Help:      "Duration of HTTP requests in seconds",
				Buckets:   prometheus.ExponentialBucketsRange(0.006099375, 16.075897965, 10),
				Namespace: namespace,
			},
			[]string{"method", "path"},
		),
	}
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(metrics.HTTPTotalNumberOfRequests)
	reg.MustRegister(metrics.HTTPRequestDuration)

	return metrics
}
