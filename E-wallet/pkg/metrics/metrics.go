package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MetricDBRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ewallet",
		Subsystem: "generic",
		Name:      "db_duration",
	}, []string{"method"})

	MetricErrCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ewallet",
		Subsystem: "generic",
		Name:      "err_count",
	}, []string{"method"})

	MetricHTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ewallet",
		Subsystem: "generic",
		Name:      "err_count",
	}, []string{"method"})
)
