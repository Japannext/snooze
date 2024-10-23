package mq

import (
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

const X_SNOOZE_STORE_INDEX = "X-Snooze-Store-Index"

var (
	storeSubOnce, storePubOnce, storeStreamOnce sync.Once
	storeStream jetstream.Stream
	storeSub *Sub
	storePub *Pub
)

func getStoreStream() jetstream.Stream {
	storeStreamOnce.Do(func() {
		storeStream = client.setupStream(jetstream.StreamConfig{
			Name: "STORE",
			Retention: jetstream.WorkQueuePolicy,
			Subjects: []string{"STORE.items"},
		})
	})
	return storeStream
}

func StoreSub() *Sub {
	storeSubOnce.Do(func() {
		stream := getStoreStream()
		storeSub = client.Consumer(stream, jetstream.ConsumerConfig{
			Name: "writer",
			Durable: "writer",
		})
	})
	return storeSub
}

func StorePub() *Pub {
	storePubOnce.Do(func() {
		getProcessStream()
		storePub = client.Pub().WithSubject("STORE.items")
	})
	return storePub
}
