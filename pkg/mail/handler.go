package mail

import (
	"bytes"
	"fmt"
	gomail "gopkg.in/mail.v2"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
)

func handler(notification *models.Notification) error {
	log.Debug("Handing notification")
	profile, found := profiles[notification.Destination.Profile]
	if !found {
		return fmt.Errorf("Dropping: profile '%s' not found", notification.Destination.Profile)
	}
	log.Debugf("Found profile %s", profile.Name)

	body, err := computeTemplate(profile, notification)
	if err != nil {
		return err
	}

	// Building email
	msg := gomail.NewMessage()
	msg.SetHeader("From", profile.From)
	msg.SetHeader("To", profile.To)
	msg.SetHeader("Subject", fmt.Sprintf("[Snooze] Alert %s", profile.Name))
	msg.SetBody("text/plain", string(body))

	tlsConfig := config.TLS.Config()
	tlsConfig.ServerName = profile.Server
	dialer := gomail.Dialer{Host: profile.Server, Port: profile.Port, TLSConfig: tlsConfig}
	sender, err := dialer.Dial()
	if err != nil {
		return err
	}

	log.Debugf("Sending mail to %s...", profile.To)
	err = sender.Send(profile.From, []string{profile.To}, msg)
	if err != nil {
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
