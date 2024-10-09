package syslog

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/utils"
)

var publishQueue = utils.NewBatchingChannel[models.Log](50, time.Second)

type Publisher struct {
	*rabbitmq.Producer
}

func NewPublisher() *Publisher {
	return &Publisher{
		rabbitmq.NewLogProducer(),
	}
}

func (pub *Publisher) Run() error {
	for batch := range publishQueue.Channel() {
		batch.TimestampMillis = uint64(time.Now().UnixMilli())
		if err := pub.Publish(batch); err != nil {
			log.Warnf("error sending batch: %s", err)
		}
		ingestedLogs.WithLabelValues(SOURCE_KIND, config.InstanceName).Add(float64(len(batch.Items)))
		batchSize.WithLabelValues(SOURCE_KIND, config.InstanceName).Observe(float64(len(batch.Items)))
	}

	return nil
}

func (pub *Publisher) Stop() {
	publishQueue.Stop()
}
