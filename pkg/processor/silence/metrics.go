package silence

import (
	"github.com/prometheus/client_golang/prometheus"
)

var silencedLogs = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "snooze",
	Name:      "silenced_logs",
	Help:      "number of logs silenced",
})

func initMetrics() {
	prometheus.MustRegister(silencedLogs)
}
