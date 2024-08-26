package rabbitmq

import (
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
}

type Handler = func(amqp.Delivery) error

type ConsumerOptions struct {
    AutoAck   bool
    Exclusive bool
    NoLocal   bool
    NoWait    bool
    Args      map[string]interface{}
}

func (consumer *Consumer) ConsumeForever() error {
	channel, err := Client.conn.Channel()
	if err != nil {
		return err
	}

	opts := consumer.Options
	if true {
		deliveries, err := channel.Consume(consumer.Queue, consumer.Topic, opts.AutoAck, opts.Exclusive, opts.NoLocal, opts.NoWait, opts.Args)
		if err != nil {
			return err
		}

		for delivery := range deliveries {
			log.Debugf("Received AMQP message: %+v", delivery.Body)
			if err := consumer.Handler(delivery); err != nil {
				log.Warnf("Rejecting message (%s): Body: %s", err, delivery.Body)
				delivery.Reject(NO_REQUEUE)
			}
		}
	}
	return nil
}
