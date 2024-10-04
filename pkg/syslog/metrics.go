package syslog

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    ingestedLogs = prometheus.NewCounterVec(prometheus.CounterOpts{
        Namespace: "snooze",
        Name: "ingested_logs",
        Help: "number of logs ingested by source plugins (and queued)",
	}, []string{"source_kind", "source_name"})
	batchSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name: "batch_size",
		Help: "size of the batch sent",
		Buckets: []float64{1.0, 5.0, 10.0, 20.0, 30.0, 50.0},
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
         Buckets: prometheus.ExponentialBuckets(0.1, 3, 8),
     }, []string{"source_kind", "source_name"})
)

func initMetrics() {
	prometheus.MustRegister(
		ingestedLogs,
		emptyTimestamp,
		sourceDelay,
		batchSize,
	)
}
