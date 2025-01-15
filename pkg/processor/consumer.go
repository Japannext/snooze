package processor

import (
	"encoding/json"

	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/models"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
)

type Consumer struct{
	manager *Manager
	processQ *mq.Sub
}

func NewConsumer() *Consumer {
	return &Consumer{
		manager: NewManager(),
		processQ: mq.ProcessSub(),
	}
}

func (c *Consumer) Run() error {
	for {
		log.Debugf("Waiting for a worker to be ready...")
		size := pool.Ready()
		log.Debugf("%d workers are ready", size)
		msgs, err := processQ.Fetch(size)
		if err != nil {
			log.Warnf("failed to fetch items: %s", err)
			continue
		}
		log.Debugf("Fetched %d logs", len(msgs))
		if len(msgs) == 0 {
			continue
		}

		for _, msgWithCtx := range msgs {
			msg := msgWithCtx.Msg
			msgCtx := msgWithCtx.Context
			// Immediately requeue messages when there is no worker
			if ok := pool.TryAcquire(); !ok {
				log.Warnf("Requeuing since no worker is ready (%d/%d)", pool.Busy(), pool.Max())
				msg.Nak()
				continue
			}
			go func() {
				defer pool.Release()

				// Get a processor with the latest config
				p := c.manager.GetProcessor()
				// Process the message
				p.Process(msgCtx, unmarshalLog(msg))
				// Ack
				msg.Ack()
			}()
		}
	}
}

func (c *Consumer) Stop() {
}

func unmarshalLog(msg jetstream.Msg) *models.Log {
	var item models.Log
	data := msg.Data()
	if err := json.Unmarshal(data, &item); err != nil {
		log.Warnf("invalid JSON while parsing log: %s", err)
		item = models.Log{
			ActualTime:  models.TimeNow(),
			DisplayTime: models.TimeNow(),
			Identity:    map[string]string{"snooze.internal": "error"},
			Message:     string(data),
			Error:       err.Error(),
		}
	}
	return &item
}
