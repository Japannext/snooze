package syslog

import (
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/mq"
)

var publishQueue = make(chan *models.Log)

type Publisher struct {
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (pub *Publisher) Run() error {
	for item := range publishQueue {
		future, err := mq.PublishAsync("PROCESS.logs", item)
		if err != nil {
			log.Warnf("failed to publish: %s", err)
			continue
		}
		go func() {
			if err := <-future.Err(); err != nil {
				log.Warnf("failed to publish (async): %s", err)
			}
		}()
		ingestedLogs.WithLabelValues(SOURCE_KIND, config.InstanceName).Inc()
	}

	return nil
}

func (pub *Publisher) Stop() {
	close(publishQueue)
}
