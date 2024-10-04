package processor

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
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
		Buckets: prometheus.ExponentialBuckets(0.1, 10, 8),
	})
	processTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name: "process_time_seconds",
		Help: "time spent processing logs",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 8),
	})
)

func processMetrics(start time.Time, item *api.Log) {
	observedTimestamp := time.UnixMilli(int64(item.ObservedTimestampMillis))

	// Process time
	processTime.Observe(time.Since(start).Seconds())

	// In-queue time
	inqueueSeconds := start.Sub(observedTimestamp).Seconds()
	inqueueTime.Observe(inqueueSeconds)

	// Counters
	processedLogs.Inc()
}

func initMetrics() {
	prometheus.MustRegister(
		processedLogs,
		inqueueTime,
		processTime,
	)
}
