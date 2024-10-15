package writer

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	writeItems = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "snooze",
		Name: "write_items",
		Help: "number of items (logs, alerts, notifications, etc) stored to the opensearch backend",
	}, []string{"kind"})
	errorItems = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "snooze",
		Name: "write_error",
		Help: "number of items that has been rejected by opensearch, and will be requeued in a simpler format",
	}, []string{"kind"})
)

func initMetrics() {
	prometheus.MustRegister(
		writeItems,
		errorItems,
	)
}
