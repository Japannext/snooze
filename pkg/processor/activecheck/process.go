package activecheck

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	log "github.com/sirupsen/logrus"
)

type Processor struct {}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	url := item.ActiveCheckURL
	if url == "" { // not an active check
		return decision.OK()
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

		return decision.Abort(err)
	}

	client := http.Client{Timeout: 1 * time.Second}

	body := bytes.NewBuffer(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return decision.Abort(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warnf("active check failed to be sent back to %s: %s", url, err)

		return decision.Abort(err)
	}
	defer resp.Body.Close()

	return decision.OK()
}
