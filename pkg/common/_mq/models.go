package mq

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/models"
)

type Batch[T models.Object] struct {
	msgs []jetstream.Msg
	items []T
}

func (batch *Batch[T]) AckAll() {
	for _, msg := range batch.msgs {
		msg.Ack()
	}
}

func (batch *Batch[T]) RescheduleAll(delay time.Duration) {
	for _, msg := range batch.msgs {
		msg.NakWithDelay(delay)
	}
}

func NewBatch[T models.Object](msgs []jetstream.Msg) *Batch[T] {
	n := len(msgs)
	batch := &Batch[T]{
		msgs: make([]jetstream.Msg, n),
		items: make([]T, n),
	}
	for _, msg := range msgs {
		var item T
		data := msg.Data()
		if err := json.Unmarshal(data, &item); err != nil {
			log.Warnf("invalid JSON while parsing item %T: %s", item, err)
			continue
		}
		batch.msgs = append(batch.msgs, msg)
		batch.items = append(batch.items, item)
	}
	return batch
}

func (batch *Batch[T]) Reject(i int) {
	if i >= len(batch.msgs) {
		return
	}
	batch.msgs[i].Term()
	removeIndex(batch.msgs, i)
	removeIndex(batch.items, i)
}

func (batch *Batch[T]) Items() []T {
	return batch.items
}

func removeIndex[T any](a []T, index int) []T {
	return append(a[:index], a[index+1:]...)
}
