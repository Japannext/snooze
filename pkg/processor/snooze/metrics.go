package snooze

import (
	"github.com/prometheus/client_golang/prometheus"
)

var snoozedLogs = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "snooze",
	Name:      "snoozed_logs",
	Help:      "number of logs snoozed",
})

func initMetrics() {
	prometheus.MustRegister(snoozedLogs)
}
