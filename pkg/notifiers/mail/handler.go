package mail

import (
	"bytes"
	"fmt"
	"net/smtp"

	log "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func computeTemplate(profile *Profile, notification *api.Notification) ([]byte, error) {
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

func handler(notification *api.Notification) error {
	profile, found := profiles[notification.Profile]
	if !found {
		return fmt.Errorf("Dropping: profile '%s' not found", notification.Profile)
	}
	addr := fmt.Sprintf("%s:%d", profile.Server, profile.Port)

	message, err := computeTemplate(profile, notification)
	if err != nil {
		return err
	}

	err = smtp.SendMail(addr, nil, profile.From, []string{profile.To}, message)
	if err != nil {
		log.Errorf("error sending mail with profile '%s': %s", profile.Name, err)
		return err
	}
	return nil
}
