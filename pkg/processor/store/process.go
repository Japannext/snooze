package store

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
)

func Process(item *api.Log) error {
	if item.Mute.SkipStorage {
		return nil
	}
	err := opensearch.LogStore.StoreLog(item)
	if err != nil {
		return err
	}
	storedLogs.Inc()
	log.Debugf("Successfully stored log %s", item)
	return nil
}

