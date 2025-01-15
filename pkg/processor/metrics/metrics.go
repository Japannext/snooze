package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

//nolint:gochecknoglobals
var (
	StoredLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snooze",
		Name:      "stored_logs_total",
		Help:      "number of logs stored in the storage (not-dropped)",
	})
	SilencedLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snooze",
		Name:      "silenced_logs_total",
		Help:      "number of logs silenced",
	})
	SnoozedLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snooze",
		Name:      "snoozed_logs_total",
		Help:      "number of logs snoozed",
	})
	ProcessedLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snooze",
		Name:      "processed_logs_total",
		Help:      "number of logs processed",
	})
	ProcessTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name:      "process_time_seconds",
		Help:      "time spent processing one log",
	})
	WorkerBusy = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "snooze",
		Name:      "worker_busy",
		Help:      "Number of workers that are currently busy",
	}, func() float64 {
		return float64(0)
	})
	WorkerMax = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "snooze",
		Name:      "worker_max",
		Help:      "Number of max workers in the pool",
	}, func() float64 {
		return float64(0)
	})
	RatelimitedLogs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "snooze",
		Name:      "ratelimited_logs_total",
		Help:      "number of logs dropped because of rate limits",
	}, []string{"name"})
)

func Init() {
	prometheus.MustRegister(
		StoredLogs,
		SilencedLogs,
		SnoozedLogs,
		ProcessedLogs,
		ProcessTime,
		WorkerBusy,
		WorkerMax,
		RatelimitedLogs,
	)
}
