package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/japannext/snooze/pkg/models"
)

const (
	nexName = "notification-v2"
)
var (
	notificationQueues = make(map[string]NotificationQueue)
)


type NotificationChannel struct {
	*amqp.Channel
}

func NewNotificationChannel() *NotificationChannel {
	var err error

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = exchange(ch, nexName, "topic", &exOpts{Durable: true})
	if err != nil {
		log.Fatal(err)
	}

	nc := &NotificationChannel{ch}
	channelsToClose["notification"] = nc
	return nc
}

type NotificationQueue struct {
	Name string
	ch   *NotificationChannel
}

func (nc *NotificationChannel) NewQueue(name string) *NotificationQueue {
	// Do not duplicate queues
	if nq, found := notificationQueues[name]; found {
		return &nq
	}

	ch := nc.Channel
	nq := NotificationQueue{name, nc}
	_, err := queue(ch, nq.Name, &qOpts{Durable: true})
	if err != nil {
		log.Fatal(err)
	}
	err = ch.QueueBind(nq.Name, "", nexName, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	notificationQueues[name] = nq
	return &nq
}

type NotificationMessage struct {
	Delivery     *amqp.Delivery
	Notification *models.Notification
}

func (nq *NotificationQueue) Consume() (<-chan NotificationMessage, error) {
	out := make(chan NotificationMessage)
	dd, err := nq.ch.Consume(nq.Name, "", true, false, false, false, nil)
	if err != nil {
		return out, err
	}
	for d := range dd {
		var notif *models.Notification
		if err := json.Unmarshal(d.Body, &notif); err != nil {
			log.Warnf("Rejecting message (%s): %s", err, d.Body)
			d.Reject(false)
		}
		out <- NotificationMessage{&d, notif}
	}
	return out, nil
}

type NotificationHandler = func(*models.Notification) error

func (ch *NotificationQueue) ConsumeForever(handler NotificationHandler) error {
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
			var item *models.Log
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

func (nq *NotificationQueue) Publish(ctx context.Context, notif *models.Notification) error {
	body, err := json.Marshal(notif)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/vnd.snooze.notification.v2+json",
		Body:         body,
	}
	return nq.ch.Publish(nexName, nq.Name, false, false, msg)
}
