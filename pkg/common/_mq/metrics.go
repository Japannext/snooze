package mq

import (
    "github.com/prometheus/client_golang/prometheus"
)

var inQueue = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "snooze",
	Name: "mq_inqueue",
	Help: "time spent in-queue",
}, []string{"stream_name"})
