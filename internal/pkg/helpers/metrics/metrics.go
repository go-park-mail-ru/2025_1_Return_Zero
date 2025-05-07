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

/*
request_duration_seconds_bucket{le="0.0003"} 73
request_duration_seconds_bucket{le="0.00039999999999999996"} 77
request_duration_seconds_bucket{le="0.0005"} 82
request_duration_seconds_bucket{le="0.00075"} 83
request_duration_seconds_bucket{le="0.001"} 84
request_duration_seconds_bucket{le="0.003"} 137
request_duration_seconds_bucket{le="0.005"} 217
request_duration_seconds_bucket{le="0.0075"} 872 ----
request_duration_seconds_bucket{le="0.01"} 1596
request_duration_seconds_bucket{le="0.030000000000000002"} 5308
request_duration_seconds_bucket{le="0.05"} 5725
request_duration_seconds_bucket{le="0.07500000000000001"} 5760
request_duration_seconds_bucket{le="0.1"} 5773
request_duration_seconds_bucket{le="0.3"} 5813
request_duration_seconds_bucket{le="0.5"} 5819
request_duration_seconds_bucket{le="0.75"} 5821
request_duration_seconds_bucket{le="1.0"} 5823
request_duration_seconds_bucket{le="+Inf"} 5827
request_duration_seconds_count 5827
request_duration_seconds_sum 112.082849507
*/

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
				Buckets:   []float64{0.0003, 0.003, 0.004, 0.005, 0.0075, 0.009, 0.01, 0.015, 0.02, 0.025, 0.03, 0.05, 1.0},
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
