package metrics

import (
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	HTTPTotalNumberOfRequests *prometheus.CounterVec
	HTTPRequestDuration       *prometheus.HistogramVec
	ActiveGoroutines          prometheus.Gauge
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
				Buckets:   prometheus.DefBuckets,
				Namespace: namespace,
			},
			[]string{"method", "path"},
		),
		ActiveGoroutines: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "active_goroutines",
                Help:      "Current number of active goroutines",
            },
        ),
	}
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(metrics.HTTPTotalNumberOfRequests)
	reg.MustRegister(metrics.HTTPRequestDuration)
	reg.MustRegister(metrics.ActiveGoroutines)

	return metrics
}