package transform

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry
var transforms []*Transform

func Startup(trs []*Transform) {
	log = logrus.WithFields(logrus.Fields{"module": "transform"})
	for _, tr := range trs {
		tr.Load()
		transforms = append(transforms, tr)
	}
}
