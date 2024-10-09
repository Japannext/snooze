package processor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	processedLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "snooze",
		Name: "processed_logs",
		Help: "number of logs processed",
	})
	inqueueTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name: "inqueue_time_seconds",
		Help: "time spent in the processing queue, waiting to be processed",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
	})
	batchTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name: "process_batch_seconds",
		Help: "time spent processing one batch (time processing batch / batch size)",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 8),
	})
	batchSize = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name: "process_batch_size",
		Help: "size of a batch",
		Buckets: []float64{1.0, 5.0, 10.0, 25.0, 50.0},
	})
)

func initMetrics() {
	prometheus.MustRegister(
		processedLogs,
		inqueueTime,
		batchTime,
		batchSize,
	)
}
