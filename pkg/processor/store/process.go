package store

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"module": "store"})

func Process(item *api.Log) error {
	err := opensearch.LogStore.StoreLog(item)
	if err != nil {
		return err
	}
	log.Debugf("Successfully stored log %s", item)
	return nil
}

