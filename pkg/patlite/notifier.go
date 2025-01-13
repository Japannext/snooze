package patlite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/japannext/snooze/pkg/common/opensearch/format"
	pl "github.com/japannext/snooze/pkg/common/patlite"
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
		return err
	}
	profileName := notification.Destination.Profile
	profile, found := profiles[profileName]
	if !found {
		return fmt.Errorf("Failed to find profile '%s'", profileName)
	}

	client, err := pl.NewClient(profile.Address, profile.Port)
	if err != nil {
		return fmt.Errorf("In profile '%s': %s", profile.Name, err)
	}
	if err := client.SetState(profile.State); err != nil {
		return fmt.Errorf("In profile '%s': %s", profile.Name, err)
	}

	notification.NotificationTime = models.TimeNow()
	err = storeQ.PublishData(ctx, &format.Create{
		Index: models.NotificationIndex,
		Item:  notification,
	})
	if err != nil {
		log.Warnf("failed to publish notification")
		tracing.Error(span, err)
		return err
	}
	return nil
}
