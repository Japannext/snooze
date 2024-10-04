package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	NO_REQUEUE = false
)

type Delivery = amqp.Delivery

type Consumer struct {
	Queue string
	Topic string
	Handler Handler
	Options ConsumerOptions

	channel *amqp.Channel
	stopping bool
	done chan bool
}

type Handler = func(amqp.Delivery) error

type ConsumerOptions struct {
    AutoAck   bool
    Exclusive bool
    NoLocal   bool
    NoWait    bool
    Args      map[string]interface{}
}

func NewConsumer(queue, topic string, handler Handler, options ConsumerOptions) *Consumer {
	channel, err := Client.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return &Consumer{
		Queue: queue,
		Topic: topic,
		Handler: handler,
		Options: options,
		channel: channel,
		done: make(chan bool),
	}
}

func (consumer *Consumer) ConsumeForever() error {
	opts := consumer.Options
	for true {
		if consumer.stopping {
			consumer.done <- true
			return nil
		}
		log.Debug("Starting channel consume loop")
		deliveries, err := consumer.channel.Consume(consumer.Queue, consumer.Topic, opts.AutoAck, opts.Exclusive, opts.NoLocal, opts.NoWait, opts.Args)
		if err != nil {
			return err
		}

		for delivery := range deliveries {
			log.Debug("Received AMQP message")
			if err := consumer.Handler(delivery); err != nil {
				log.Warnf("Rejecting message (%s): Body: %s", err, delivery.Body)
				delivery.Reject(NO_REQUEUE)
			}
		}
		log.Debug("Done handling deliveries")
	}
	log.Debug("Exited channel consume loop")
	return nil
}

func (consumer *Consumer) Name() string {
	return "consumer"
}

func (consumer *Consumer) Pass() error {
	if consumer.channel.IsClosed() {
		return fmt.Errorf("channel closed")
	}
	return nil
}

func (consumer *Consumer) Stop() {
	log.Debugf("Stopping consumer...")
	consumer.stopping = true
	<-consumer.done
	log.Debugf("Stopped consumer")
}
