package store

import (
	"bytes"
	"encoding/json"
	"fmt"

	api "github.com/japannext/snooze/common/api/v2"
	"github.com/japannext/snooze/common/opensearch"
)

const index = "alerts-v2"

var client *opensearch.OpensearchClient

func Init() {
	client = opensearch.Client
}

func Process(alert *api.Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	resp, err := client.Index(index, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("Unexpected HTTP status code: %s", resp.String())
	}
	return nil
}
