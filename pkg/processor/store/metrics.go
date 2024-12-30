package store

import (
	"github.com/prometheus/client_golang/prometheus"
)

var storedLogs = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "snooze",
	Name:      "stored_logs",
	Help:      "number of logs stored in the storage (not-dropped)",
})

func initMetrics() {
	prometheus.MustRegister(storedLogs)
}
