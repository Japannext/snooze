package googlechat

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/common/opensearch"
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
	if err := storeQ.PublishData(ctx, opensearch.Create(models.NOTIFICATION_INDEX, notification)); err != nil {
		log.Warnf("failed to publish notification")
		tracing.Error(span, err)
		return err
	}
	return nil
}
