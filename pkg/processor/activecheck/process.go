package activecheck

import (
	"context"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(ctx context.Context, item *api.Log) error {

	url, ok := item.Labels["activecheck.url"]
	if !ok {
		return nil
	}

	// Making sure the log is dropped no matter what
	item.Mute.Drop("active check")

	delayMillis := uint64(time.Now().UnixMilli()) - item.TimestampMillis
	data, err := json.Marshal(api.ActiveCheckCallback{
		DelayMillis: delayMillis,
		Error: item.Error,
	})
	if err != nil {
		log.Warnf("failed to marshal response: %s", err)
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Warnf("active check failed to be sent back to %s: %s", url, err)
		return err
	}

	return nil
}
