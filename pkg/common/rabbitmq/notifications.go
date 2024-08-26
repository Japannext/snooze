package rabbitmq

const (
	NOTIFICATION_EXCHANGE = "notification-v2"
	NOTIFICATION_CONTENT_TYPE = "application/vnd.snooze.notification.v2+json"
)


func SetupNotifications(queueNames []string) {
	// Setup
	exchanges := map[string]ExchangeOptions{
		NOTIFICATION_EXCHANGE: {Kind: "topic", Durable: true},
	}
	queues := map[string]QueueOptions{}
	for _, queueName := range queueNames {
		queues[queueName] = QueueOptions{Durable: true}
	}
	binds := map[string]string{}
	for _, queueName := range queueNames {
		binds[queueName] = NOTIFICATION_EXCHANGE
	}
	Client.setup(exchanges, queues, binds)
}

func NewNotificationProducer(queue string) *Producer {
	return NewProducer(NOTIFICATION_EXCHANGE, queue, "", NOTIFICATION_CONTENT_TYPE)
}
