package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewMockMetrics() *Metrics {
	return &Metrics{
		GRPCTotalNumberOfRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "mock_requests"},
			[]string{"method", "status"},
		),
		GRPCRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "mock_duration"},
			[]string{"method"},
		),
		DatabaseDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "mock_db_duration"},
			[]string{"operation"},
		),
		DatabaseErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "mock_db_errors"},
			[]string{"operation"},
		),
	}
}
