package googlechat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/japannext/snooze/pkg/common/opensearch/format"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
)

func notificationHandler(ctx context.Context, msg jetstream.Msg) error {
	ctx, span := tracer.Start(ctx, "handler")
	defer span.End()

	var notification *models.Notification
	if err := json.Unmarshal(msg.Data(), &notification); err != nil {
		tracing.Error(span, err)
		return err
	}
	profileName := notification.Destination.Profile
	profile, found := profiles[profileName]
	if !found {
		return fmt.Errorf("Failed to find profile '%s'", profileName)
	}

	chatMsg := profile.FormatToMessage(notification)
	if err := client.SendMessage(ctx, profile.Space, chatMsg); err != nil {
		log.Warnf("error sending message: %s", err)
		tracing.Error(span, err)
		return err
	}

	notification.NotificationTime = models.TimeNow()
	err := storeQ.PublishData(ctx, &format.Create{
		Index: models.NOTIFICATION_INDEX,
		Item:  notification,
	})
	if err != nil {
		log.Warnf("failed to publish notification")
		tracing.Error(span, err)
		return err
	}
	return nil
}
