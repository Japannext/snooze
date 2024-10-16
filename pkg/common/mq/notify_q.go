package mq

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

var (
	notifySubOnce, notifyPubOnce, notifyStreamOnce sync.Once
	notifyStream jetstream.Stream
	notifySub *Sub
	notifyPub *Pub
)

func getNotifyStream() jetstream.Stream {
	notifyStreamOnce.Do(func() {
		client := getClient()
		subjects := []string{}
		for _, backend := range strings.Split(config.NotifyBackends, ",") {
			subjects = append(subjects, fmt.Sprintf("NOTIFY.%s", backend))
		}
		client.setupStream(jetstream.StreamConfig{
			Name: "NOTIFY",
			Retention: jetstream.WorkQueuePolicy,
			Subjects: subjects,
		})
	})
	return notifyStream
}

func NotifySub(backend string) *Sub {
	notifySubOnce.Do(func() {
		client := getClient()
		stream := getNotifyStream()
		notifySub = client.Subscribe(stream, backend)
	})
	return notifySub
}

func NotifyPub() *Pub {
	notifyPubOnce.Do(func() {
		client := getClient()
		getNotifyStream()
		notifyPub = client.Pub("PROCESS.logs")
	})
	return notifyPub
}
