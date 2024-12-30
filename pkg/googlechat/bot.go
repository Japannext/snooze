package googlechat

// A googlechat bot, which is a daemon listening
// to a pubsub, reacting to events

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	log "github.com/sirupsen/logrus"
	chat "google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

type Bot struct {
	pubsub *pubsub.Client
}

func NewBot() *Bot {
	ctx := context.Background()
	creds := getCredentials()
	pubsubClient, err := pubsub.NewClient(ctx, "snooze-v2", option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("failed to initialize pubsub: %s", err)
	}
	return &Bot{
		pubsub: pubsubClient,
	}
}

func (bot *Bot) Run() error {

	ctx := context.Background()

	// topic := bot.pubsub.Topic("snooze")
	sub := bot.pubsub.Subscription("snooze")

	err := sub.Receive(ctx, pubsubRouter)
	if err != nil {
		return err
	}

	return nil
}

// Route the pubsub event to the correct function
func pubsubRouter(ctx context.Context, msg *pubsub.Message) {

	var event *chat.DeprecatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		// TODO
	}

	switch event.Type {
	case "CARD_CLICKED":
		// TODO
	case "SUBMIT_FORM":
		// TODO
	default:
		// TODO
	}

	msg.Ack()
}

func onCardClick(event *chat.DeprecatedEvent) {
	switch event.Common.InvokedFunction {
	case "openSnoozeDialog":
		openSnoozeDialog(event)
	}
}

func openSnoozeDialog(event *chat.DeprecatedEvent) {
}

func (bot *Bot) Stop() {
}
