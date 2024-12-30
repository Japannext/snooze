package ratelimit

import (
	"github.com/prometheus/client_golang/prometheus"
)

var rateLimitedLogs = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "snooze",
	Name:      "ratelimited_logs",
	Help:      "number of logs dropped because of rate limits",
}, []string{"name"})

func initMetrics() {
	prometheus.MustRegister(rateLimitedLogs)
}
