package silence

import (
	"github.com/sirupsen/logrus"
)

var silences []*Silence
var log *logrus.Entry

func Startup(rules []*Silence) {
	initMetrics()
	log = logrus.WithFields(logrus.Fields{"module": "silence"})
	for _, s := range rules {
		s.Load()
		silences = append(silences, s)
	}
}
