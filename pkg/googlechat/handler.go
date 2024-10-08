package googlechat

import (
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/japannext/snooze/pkg/models"
)

func handler(notification *models.Notification) error {
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
