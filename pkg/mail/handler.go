package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/japannext/snooze/pkg/models"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

func handler(ctx context.Context, msg jetstream.Msg) error {
	var notification *models.Notification
	if err := json.Unmarshal(msg.Data(), &notification); err != nil {
		return err
	}

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
	mail := gomail.NewMessage()
	mail.SetHeader("From", profile.From)
	mail.SetHeader("To", profile.To)
	mail.SetHeader("Subject", fmt.Sprintf("[Snooze] Alert %s", profile.Name))
	mail.SetBody("text/plain", string(body))

	tlsConfig := config.TLS.Config()
	tlsConfig.ServerName = profile.Server
	dialer := gomail.Dialer{Host: profile.Server, Port: profile.Port, TLSConfig: tlsConfig}
	sender, err := dialer.Dial()
	if err != nil {
		return err
	}

	log.Debugf("Sending mail to %s...", profile.To)
	err = sender.Send(profile.From, []string{profile.To}, mail)
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
