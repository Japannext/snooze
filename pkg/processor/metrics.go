package processor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	processedLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snooze",
		Name:      "processed_logs",
		Help:      "number of logs processed",
	})
	processTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name:      "process_time_seconds",
		Help:      "time spent processing one log",
	})
	workerBusy = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "snooze",
		Name:      "worker_busy",
		Help:      "Number of workers that are currently busy",
	}, func() float64 {
		return float64(pool.Busy())
	})
	workerMax = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "snooze",
		Name:      "worker_max",
		Help:      "Number of max workers in the pool",
	}, func() float64 {
		return float64(pool.Max())
	})
)

func initMetrics() {
	prometheus.MustRegister(
		processedLogs,
		processTime,
		workerBusy,
		workerMax,
	)
}
