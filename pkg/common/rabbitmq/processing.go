package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const (
	pexName = "processing-log-v2"
	pqName  = "processing-log-v2"
)

type ProcessChannel struct {
	*amqp.Channel
	done chan struct{}
	stopping bool
}

func InitProcessChannel() *ProcessChannel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	err = exchange(ch, pexName, "direct", &exOpts{Durable: true})
	if err != nil {
		log.Fatal(err)
	}
	_, err = queue(ch, pqName, &qOpts{Durable: true})
	if err != nil {
		log.Fatal(err)
	}
	err = ch.QueueBind(pqName, "", pexName, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	pch := &ProcessChannel{Channel: ch, done: done, stopping: false}
	channelsToClose["process"] = pch
	return pch
}

func (ch *ProcessChannel) Cancel() error {
	return ch.Channel.Cancel(pqName, false)
}

type LogMessage struct {
	Delivery *amqp.Delivery
	Log    *api.Log
}

type LogHandler = func(*api.Log) error

func (ch *ProcessChannel) ConsumeForever(handler LogHandler) error {
	for true {
		if ch.stopping {
			log.Debugf("Stopping ConsumeForever loop")
			break
		}
		deliveries, err := consume(ch.Channel, pqName, "", &chOpts{AutoAck: true})
		if err != nil {
			return err
		}

		for d := range deliveries {
			if ch.stopping {
				log.Debugf("Stopping delivery channel loop")
				break
			}
			log.Debugf("Received an AMQP message!")
			var item *api.Log
			if err := json.Unmarshal(d.Body, &item); err != nil {
				log.Warnf("Rejecting message (%s): %s", err, d.Body)
				// discard(d)
			}
			log.Debugf("Appending log to channel")
			if err := handler(item); err != nil {
				log.Errorf("error handling message: %s", err)
				// discard(d)
			}
			log.Debug("Done processing message")
		}
		log.Debug("Exited delivery loop")
	}
	log.Debugf("Sending termination signal")
	ch.done <- struct{}{}
	return nil
}

/*
func (ch *ProcessChannel) Consume() (<-chan LogMessage, error) {
	log.Debugf("Consuming!")
	out := make(chan LogMessage)
	dd, err := consume(ch.Channel, pqName, "", &chOpts{AutoAck: true})
	if err != nil {
		return out, err
	}
	for d := range dd {
		log.Debugf("Received an AMQP message!")
		var item *api.Log
		if err := json.Unmarshal(d.Body, &item); err != nil {
			log.Warnf("Rejecting message (%s): %s", err, d.Body)
			d.Reject(false)
		}
		log.Debugf("Appending log to channel")
		out <- LogMessage{&d, item}
	}
	log.Warnf("Out of consume loop!")
	return out, nil
}
*/

func (ch *ProcessChannel) Publish(ctx context.Context, item *api.Log) error {
	body, err := json.Marshal(item)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/vnd.snooze.log.v2+json",
		Body:         body,
	}
	return ch.Channel.Publish(pqName, "", false, false, msg)
}

func (ch *ProcessChannel) Stop() {
	ch.stopping = true
	ch.Cancel()
	<- ch.done
}
