package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var Client = &client{}

type client struct {
	conn *amqp.Connection
}

type ExchangeOptions struct {
	Kind       string
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

type BindOptions struct {
	Exchange string
	Key string
	NoWait bool
	Args map[string]interface{}
}

func (c *client) setup(exchanges map[string]ExchangeOptions, queues map[string]QueueOptions, binds map[string]BindOptions) {
	channel, err := c.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	for name, cfg := range exchanges {
		log.Debugf("Ensuring exchange '%s'", name)
		err := channel.ExchangeDeclare(name, cfg.Kind, cfg.Durable, cfg.AutoDelete, cfg.Internal, cfg.NoWait, cfg.Args)
		if err != nil {
			log.Fatal(err)
		}
	}
	for name, cfg := range queues {
		log.Debugf("Ensuring queue '%s'", name)
		_, err := channel.QueueDeclare(name, cfg.Durable, cfg.DeleteWhenUnused, cfg.Exclusive, cfg.NoWait, cfg.Args)
		if err != nil {
			log.Fatal(err)
		}
	}
	for queue, cfg := range binds {
		log.Debugf("Ensuring bind exchange[%s] -[%s]-> queue[%s]", cfg.Exchange, cfg.Key, queue)
		if err := channel.QueueBind(queue, cfg.Key, cfg.Exchange, cfg.NoWait, cfg.Args); err != nil {
			log.Fatal(err)
		}
	}
}
