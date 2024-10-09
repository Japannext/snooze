package activecheck

import (
	"context"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "activecheck")
	defer span.End()
	url, ok := item.Labels["activecheck.url"]
	if !ok {
		return nil
	}

	// Making sure the log is dropped no matter what
	item.Mute.Drop("active check")

	delayMillis := uint64(time.Now().UnixMilli()) - item.Timestamp.Actual

	if delayMillis > uint64((30 * time.Second).Milliseconds()) {
		// skipping
		return nil
	}

	data, err := json.Marshal(models.ActiveCheckCallback{
		DelayMillis: delayMillis,
		Error: item.Error,
	})
	if err != nil {
		log.Warnf("failed to marshal response: %s", err)
		return err
	}
	go func() {
		if _, err := http.Post(url, "application/json", bytes.NewBuffer(data)); err != nil {
			log.Warnf("active check failed to be sent back to %s: %s", url, err)
		}
	}()

	return nil
}
