package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	GRPCTotalNumberOfRequests *prometheus.CounterVec
	GRPCRequestDuration       *prometheus.HistogramVec
	DatabaseDuration          *prometheus.HistogramVec
	DatabaseErrors            *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer, namespace string) *Metrics {
	metrics := &Metrics{
		GRPCTotalNumberOfRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      "grpc_total_number_of_requests",
				Help:      "Total number of requests received",
				Namespace: namespace,
			},
			[]string{"method", "status"},
		),
		GRPCRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:      "grpc_request_duration_seconds",
				Help:      "Duration of gRPC requests in seconds",
				Namespace: namespace,
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method"},
		),
		DatabaseDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:      "database_duration_seconds",
				Help:      "Duration of database operations in seconds",
				Namespace: namespace,
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		DatabaseErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      "database_errors_total",
				Help:      "Total number of database errors",
				Namespace: namespace,
			},
			[]string{"operation"},
		),
	}
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(metrics.GRPCTotalNumberOfRequests)
	reg.MustRegister(metrics.GRPCRequestDuration)
	reg.MustRegister(metrics.DatabaseDuration)
	reg.MustRegister(metrics.DatabaseErrors)

	return metrics
}
