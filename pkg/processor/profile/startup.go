package profile

import (
	"github.com/sirupsen/logrus"
)

var fastMapper *FastMapper
var log *logrus.Entry

func Startup(prfs []*Profile) {
	log = logrus.WithFields(logrus.Fields{"module": "profile"})
	fastMapper = NewFastMapper(prfs)
}
