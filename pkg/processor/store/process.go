package store

import (
	"bytes"
	"encoding/json"
	"fmt"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
)

const ALERTS = "alerts-v2"

func Process(alert *api.Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	resp, err := opensearch.Client.Index(ALERTS, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("Unexpected HTTP status code: %s", resp.String())
	}
	return nil
}
