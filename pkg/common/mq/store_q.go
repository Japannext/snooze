package mq

import (
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

var (
	storeSubOnce, storePubOnce, storeStreamOnce sync.Once
	storeStream jetstream.Stream
	storeSub *Sub
	storePub *Pub
)

func getStoreStream() jetstream.Stream {
	storeStreamOnce.Do(func() {
		client := getClient()
		client.setupStream(jetstream.StreamConfig{
			Name: "STORE",
			Retention: jetstream.WorkQueuePolicy,
			Subjects: []string{"STORE.*"},
		})
	})
	return storeStream
}

func StoreSub() *Sub {
	storeSubOnce.Do(func() {
		client := getClient()
		stream := getProcessStream()
		storeSub = client.Subscribe(stream, "writer")
	})
	return storeSub
}

func StorePub(subject string) *Pub {
	storePubOnce.Do(func() {
		client := getClient()
		getProcessStream()
		storePub = client.Pub(subject)
	})
	return storePub
}
