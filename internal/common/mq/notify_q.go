package mq

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

var (
	notifySubOnce, notifyPubOnce, notifyStreamOnce sync.Once
	notifyStream                                   jetstream.Stream
	notifySub                                      *Sub
	notifyPub                                      *Pub
)

func getNotifyStream() jetstream.Stream {
	notifyStreamOnce.Do(func() {
		notifyStream = client.setupStream(jetstream.StreamConfig{
			Name:      "NOTIFY",
			Retention: jetstream.WorkQueuePolicy,
			Subjects:  []string{"NOTIFY.*"},
			Replicas:  config.Replicas,
		})
	})
	return notifyStream
}

func NotifySub(name string) *Sub {
	notifySubOnce.Do(func() {
		stream := getNotifyStream()
		notifySub = client.Consumer(stream, jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			FilterSubject: fmt.Sprintf("NOTIFY.%s", name),
		})
	})
	return notifySub
}

func NotifyPub() *Pub {
	notifyPubOnce.Do(func() {
		getNotifyStream()
		notifyPub = client.Pub()
	})
	return notifyPub
}
