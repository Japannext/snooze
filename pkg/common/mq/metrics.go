package mq

import (
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
    "github.com/prometheus/client_golang/prometheus"
)

const X_SNOOZE_PUBLISHED_TIME = "X-Snooze-Published-Time"

var inQueue = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "snooze",
	Name: "mq_inqueue",
	Help: "time spent in-queue",
}, []string{"stream_name"})

// Inject the publication time into the message header
func injectPublishTime(msg *nats.Msg) {
	header := nats.Header(map[string][]string{
		X_SNOOZE_PUBLISHED_TIME: []string{strconv.Itoa(int(time.Now().UnixMilli()))},
	})
	msg.Header = header
}

// Use the publication time to deduce the amount of time
// the message spent in queue
func observeDelay(streamName string, msg jetstream.Msg) {
	values, ok := msg.Headers()[X_SNOOZE_PUBLISHED_TIME]
	if !ok {
		return
	}
	if len(values) == 0 {
		return
	}
	i, err := strconv.Atoi(values[0])
	if err != nil {
		return
	}
	start := time.UnixMilli(int64(i))
	delay := time.Since(start)
	inQueue.WithLabelValues(streamName).Observe(float64(delay.Seconds()))
}
