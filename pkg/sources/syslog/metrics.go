package syslog

import (
    "time"

    "github.com/prometheus/client_golang/prometheus"

    api "github.com/japannext/snooze/pkg/common/api/v2"
)

var buckets = []float64{0.1, 0.5, 1.0, 10.0, 60.0, 3600.0}

var (
    processedLogs = prometheus.NewCounterVec(prometheus.CounterOpts{
        Namespace: "snooze",
        Name: "source_processed_logs",
        Help: "number of logs processed",
	}, []string{"source_kind", "source_name"})
	emptyTimestamp = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "snooze",
		Name: "log_empty_timestamp",
		Help: "number of logs which has their timestamp set to empty value",
	}, []string{"source_kind", "source_name"})
   sourceDelay = prometheus.NewHistogramVec(prometheus.HistogramOpts{
         Namespace: "snooze",
         Name: "delay_by_source_seconds",
         Help: "time between log timestamp and entering the earlier snooze sub-system",
         Buckets: buckets,
     }, []string{"source_kind", "source_name"})
)

func processMetrics(start time.Time, item *api.Log) {
	timestamp := time.UnixMilli(int64(item.TimestampMillis))
	observedTimestamp := time.UnixMilli(int64(item.ObservedTimestampMillis))

	// counter
	processedLogs.WithLabelValues(SOURCE_KIND, config.InstanceName).Inc()

	if item.TimestampMillis == item.ObservedTimestampMillis {
		emptyTimestamp.WithLabelValues(SOURCE_KIND, config.InstanceName).Inc()
	} else {
		// Only observe delay when the timestamp is known
		delay := observedTimestamp.Sub(timestamp).Seconds()
		sourceDelay.WithLabelValues(SOURCE_KIND, config.InstanceName).Observe(delay)
	}
}

func initMetrics() {
	prometheus.MustRegister(
		processedLogs,
		emptyTimestamp,
		sourceDelay,
	)
}
