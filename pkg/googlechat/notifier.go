package googlechat

import (
	"bytes"
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
	/*
	text, err := computeTemplate(profile, notification)
	if err != nil {
		log.Warnf("error computing template: %s", err)
		tracing.Error(span, err)
		return err
	}
	*/

	msgCard := GetCard(notification)
	if err := client.SendMessage(ctx, profile.Space, msgCard); err != nil {
		log.Warnf("error sending message: %s", err)
		tracing.Error(span, err)
		return err
	}

	if err := storeQ.PublishData(ctx, opensearch.Create(models.NOTIFICATION_INDEX, notification)); err != nil {
		log.Warnf("failed to publish notification")
		tracing.Error(span, err)
		return err
	}
	return nil
}

func computeTemplate(profile *Profile, notification *models.Notification) ([]byte, error) {
    var buf bytes.Buffer
    err := profile.internal.template.Execute(&buf, notification)
    if err != nil {
        if profile.internal.isDefault {
            return []byte(""), err
        } else {
            log.Warnf("failed to execute template for profile '%s'. Fall-back to default template", profile.Name)
            buf = bytes.Buffer{}
            if err := defaultTemplate.Execute(&buf, notification); err != nil {
                return []byte(""), err
            }
        }
    }
    return buf.Bytes(), nil
}
