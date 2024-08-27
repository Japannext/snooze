package rabbitmq

const (
	PROCESSING_EXCHANGE = "processing-log-v2"
	PROCESSING_QUEUE = "processing-log-v2"
	LOG_CONTENT_TYPE = "application/vnd.snooze.log.v2+json"
)

func SetupProcessing() {
	// Setup
	exchanges := map[string]ExchangeOptions{
		PROCESSING_EXCHANGE: {Kind: "direct", Durable: true},
	}
	queues := map[string]QueueOptions{
		PROCESSING_QUEUE: {Durable: true},
	}
	binds := map[string]BindOptions{
		PROCESSING_QUEUE: {Exchange: PROCESSING_EXCHANGE},
	}
	Client.setup(exchanges, queues, binds)
}

func NewLogProducer() *Producer {
	return NewProducer(PROCESSING_EXCHANGE, "", LOG_CONTENT_TYPE)
}
