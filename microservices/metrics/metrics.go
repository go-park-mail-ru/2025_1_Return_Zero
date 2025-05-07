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

/*
request_duration_seconds_bucket{le="0.0003"} 1153
request_duration_seconds_bucket{le="0.00039999999999999996"} 1373
request_duration_seconds_bucket{le="0.0005"} 1747
request_duration_seconds_bucket{le="0.00075"} 4031
request_duration_seconds_bucket{le="0.001"} 7585
request_duration_seconds_bucket{le="0.003"} 38531
request_duration_seconds_bucket{le="0.005"} 47689
request_duration_seconds_bucket{le="0.0075"} 50412
request_duration_seconds_bucket{le="0.01"} 51028
request_duration_seconds_bucket{le="0.030000000000000002"} 51701
request_duration_seconds_bucket{le="0.05"} 51743
request_duration_seconds_bucket{le="0.07500000000000001"} 51753
request_duration_seconds_bucket{le="0.1"} 51763
request_duration_seconds_bucket{le="0.3"} 51835
request_duration_seconds_bucket{le="0.5"} 51850
request_duration_seconds_bucket{le="0.75"} 51852
request_duration_seconds_bucket{le="1.0"} 51854
request_duration_seconds_bucket{le="+Inf"} 51862
request_duration_seconds_count 51862
request_duration_seconds_sum 162.932875708
*/

func NewMetrics(reg prometheus.Registerer, namespace string) *Metrics {
	metrics := &Metrics{
		GRPCTotalNumberOfRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      "grpc_number_of_requests_total",
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
				Buckets:   []float64{0.0003, 0.0005, 0.00075, 0.001, 0.002, 0.003, 0.004, 0.005, 0.0075, 0.01, 0.3, 1},
			},
			[]string{"method"},
		),
		DatabaseDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:      "database_duration_seconds",
				Help:      "Duration of database operations in seconds",
				Namespace: namespace,
				Buckets:   []float64{0.0003, 0.0005, 0.00075, 0.001, 0.002, 0.003, 0.004, 0.005, 0.0075, 0.01, 0.3, 1},
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
