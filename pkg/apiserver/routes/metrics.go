package routes

import (
	"github.com/prometheus/client_golang/prometheus"
)

var buckets = []float64{0.1, 0.5, 1.0, 10.0}

var logSearchDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Namespace: "snooze",
	Name:      "log_search_duration_seconds",
	Help:      "time spent for each log search",
	Buckets:   buckets,
})

func InitMetrics() {
	prometheus.MustRegister(
		logSearchDuration,
	)
}
