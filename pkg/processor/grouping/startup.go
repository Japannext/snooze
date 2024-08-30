package grouping

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry
var groupings []*Grouping

func Startup(rules []*Grouping) {
	log = logrus.WithFields(logrus.Fields{"module": "grouping"})
	for _, group := range rules {
		group.Load()
		groupings = append(groupings, group)
	}
}
