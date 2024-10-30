package mq

import (
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
    "github.com/prometheus/client_golang/prometheus"
)

const X_SNOOZE_PUBLISHED_TIME = "X-Snooze-Published-Time"

var (
	inQueue = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "snooze",
		Name: "mq_inqueue",
		Help: "time spent in-queue",
	}, []string{"stream_name"})
)

// Inject the publication time into the message header
func injectPublishTime(header *nats.Header) {
	header.Add(X_SNOOZE_PUBLISHED_TIME, strconv.Itoa(int(time.Now().UnixMilli())))
}

func getPublishedTime(msg jetstream.Msg) (time.Time, bool) {
	value := msg.Headers().Get(X_SNOOZE_PUBLISHED_TIME)
	if value == "" {
		log.Warnf("header X-Snooze-Published-Time is empty")
		return time.Time{}, false
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Warnf("failed to convert `%s` to integer", value)
		return time.Time{}, false
	}
	return time.UnixMilli(int64(i)), true
}

// Use the publication time to deduce the amount of time
// the message spent in queue
func observeDelay(streamName string, start time.Time) {
	delay := time.Since(start)
	inQueue.WithLabelValues(streamName).Observe(float64(delay.Seconds()))
}

func initMetrics() {
	prometheus.MustRegister(
		inQueue,
	)
}
