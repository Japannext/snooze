package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	api "github.com/japannext/snooze/common/api/v2"
)

const (
	pexName = "processing-v1"
	pqName  = "processing-v1"
)

type ProcessChannel struct {
	*amqp.Channel
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

	pch := &ProcessChannel{ch}
	channelsToClose["process"] = pch
	return pch
}

func (ch *ProcessChannel) Cancel() error {
	return ch.Channel.Cancel(pqName, false)
}

type AlertMessage struct {
	Delivery *amqp.Delivery
	Alert    *api.Alert
}

func (ch *ProcessChannel) Consume() (<-chan AlertMessage, error) {
	out := make(chan AlertMessage)
	dd, err := consume(ch.Channel, pqName, "", &chOpts{AutoAck: true})
	if err != nil {
		return out, err
	}
	for d := range dd {
		var alert *api.Alert
		if err := json.Unmarshal(d.Body, &alert); err != nil {
			log.Warnf("Rejecting message (%s): %s", err, d.Body)
			d.Reject(false)
		}
		out <- AlertMessage{&d, alert}
	}
	return out, nil
}

func (ch *ProcessChannel) Publish(ctx context.Context, alert *api.Alert) error {
	body, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/vnd.snooze.alert.v2+json",
		Body:         body,
	}
	return ch.Channel.Publish(pqName, "", false, false, msg)
}
