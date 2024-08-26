package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var Client = &client{}

type client struct {
	conn *amqp.Connection
}

type ExchangeOptions struct {
	Kind	   string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       map[string]interface{}
}

type QueueOptions struct {
    Durable          bool
    DeleteWhenUnused bool
    Exclusive        bool
    NoWait           bool
    Args             map[string]interface{}
}

func (c *client) setup(exchanges map[string]ExchangeOptions, queues map[string]QueueOptions, binds map[string]string) {
	channel, err := c.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	for name, cfg := range exchanges {
		err := channel.ExchangeDeclare(name, cfg.Kind, cfg.Durable, cfg.AutoDelete, cfg.Internal, cfg.NoWait, cfg.Args)
		if err != nil {
			log.Fatal(err)
		}
	}
	for name, cfg := range queues {
		_, err := channel.QueueDeclare(name, cfg.Durable, cfg.DeleteWhenUnused, cfg.Exclusive, cfg.NoWait, cfg.Args)
		if err != nil {
			log.Fatal(err)
		}
	}
	for ex, queue := range binds {
		if err := channel.QueueBind(queue, "", ex, false, nil); err != nil {
			log.Fatal(err)
		}
	}
}
