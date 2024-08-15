package opensearch

import (
	"github.com/sirupsen/logrus"
)

var (
	LogStore *OpensearchLogStore
)

var log *logrus.Entry

func Init() {

	log = logrus.WithFields(logrus.Fields{"module": "opensearch"})

	LogStore = NewLogStore()

	// Fail immediately if the database is unreachable
	if err := LogStore.CheckHealth(); err != nil {
		log.Fatal(err)
	}

	if err := LogStore.Bootstrap(); err != nil {
		log.Fatal(err)
	}
}
