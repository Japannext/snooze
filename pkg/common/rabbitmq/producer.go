package rabbitmq

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp.Channel
	Exchange string
	Key string
	Mandatory bool
	Immediate bool
	ContentType string
}

func NewProducer(exchange, queue, topic, contentType string) *Producer {
	channel, err := Client.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return &Producer{
		channel: channel,
		Exchange: queue,
		Key: topic,
		ContentType: contentType,
	}
}

func (p *Producer) Publish(item interface{}) error {
	body, err := json.Marshal(item)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  p.ContentType,
		Body:         body,
	}
	return p.channel.Publish(p.Exchange, p.Key, p.Mandatory, p.Immediate, msg)
}
