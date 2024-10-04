package alertmanager

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ingestedAlerts = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "snooze",
		Name: "ingested_alerts",
		Help: "number of alerts ingested by source plugins (and queued)",
	}, []string{"source_kind", "source_name"})
	updatedAlerts = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "snooze",
		Name: "updated_alerts",
		Help: "number of alerts refreshed (exists in snooze, but still firing)",
	}, []string{"source_kind", "source_name"})
)

func initMetrics() {
	prometheus.MustRegister(
		ingestedAlerts,
		updatedAlerts,
	)
}
