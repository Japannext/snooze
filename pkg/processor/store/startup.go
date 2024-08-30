package store

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func Startup() {
	log = logrus.WithFields(logrus.Fields{"module": "store"})
	initMetrics()
}
