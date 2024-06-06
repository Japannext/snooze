package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	api "github.com/japannext/snooze/pkg/common/api/v2"
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

type AlertHandler = func(*api.Alert) error

func handleAlerts(handler AlertHandler, deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Debugf("Received an AMQP message!")
		var alert *api.Alert
		if err := json.Unmarshal(d.Body, &alert); err != nil {
			log.Warnf("Rejecting message (%s): %s", err, d.Body)
			d.Reject(false)
		}
		log.Debugf("Appending alert to channel")
		if err := handler(alert); err != nil {
			log.Errorf("error handling message: %s", err)
			d.Reject(false)
		}
	}
	log.Warnf("DONE HANDLING ALERTS")
}

func (ch *ProcessChannel) ConsumeForever(handler AlertHandler) error {
	deliveries, err := consume(ch.Channel, pqName, "", &chOpts{AutoAck: true})
	if err != nil {
		return err
	}

	var done chan error
	go handleAlerts(handler, deliveries, done)

	<-done
	return nil
}

func (ch *ProcessChannel) Consume() (<-chan AlertMessage, error) {
	log.Debugf("Consuming!")
	out := make(chan AlertMessage)
	dd, err := consume(ch.Channel, pqName, "", &chOpts{AutoAck: true})
	if err != nil {
		return out, err
	}
	for d := range dd {
		log.Debugf("Received an AMQP message!")
		var alert *api.Alert
		if err := json.Unmarshal(d.Body, &alert); err != nil {
			log.Warnf("Rejecting message (%s): %s", err, d.Body)
			d.Reject(false)
		}
		log.Debugf("Appending alert to channel")
		out <- AlertMessage{&d, alert}
	}
	log.Warnf("Out of consume loop!")
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
