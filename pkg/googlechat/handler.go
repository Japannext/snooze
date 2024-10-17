package googlechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/japannext/snooze/pkg/models"
	"github.com/nats-io/nats.go/jetstream"
)

func handler(ctx context.Context, msg jetstream.Msg) error {

    var notification *models.Notification
    if err := json.Unmarshal(msg.Data(), &notification); err != nil {
        return err
    }
	profileName := notification.Destination.Profile
	profile, found := profiles[profileName]
	if !found {
		return fmt.Errorf("Failed to find profile '%s'", profileName)
	}
	text, err := computeTemplate(profile, notification)
	if err != nil {
		log.Warnf("error computing template: %s", err)
		return err
	}
	if err := client.SendNewMessage(profile.Space, string(text)); err != nil {
		log.Warnf("error sending message: %s", err)
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
