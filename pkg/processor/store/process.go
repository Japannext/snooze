package store

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
)

const ALERTS = "alerts-v2"

func Process(alert *api.Alert) error {
	err := opensearch.LogStore.Store(alert)
	if err != nil {
		return err
	}
	return nil
}
