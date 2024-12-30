package activecheck

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) error {
	url := item.ActiveCheckURL
	if url == "" {
		return nil
	}

	// Making sure the log is dropped no matter what
	item.Status.Kind = models.LogActiveCheck
	item.Status.SkipNotification = true
	item.Status.SkipStorage = true

	data, err := json.Marshal(models.SourceActiveCheck{
		Error: item.Error,
	})
	if err != nil {
		log.Warnf("failed to marshal response: %s", err)
		return err
	}
	client := http.Client{Timeout: 1 * time.Second}
	if _, err := client.Post(url, "application/json", bytes.NewBuffer(data)); err != nil {
		log.Warnf("active check failed to be sent back to %s: %s", url, err)
	}

	return nil
}
