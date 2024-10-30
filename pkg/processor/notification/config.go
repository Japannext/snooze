package notification

import (
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
)

type Notification struct {
	Name string	`yaml:"name"`
	If string `yaml:"if"`
	Destinations []models.Destination `yaml:"destinations"`

	internal struct {
		condition *lang.Condition
	}
}

func (notif *Notification) Load() {
	var err error
	log.Debugf("Loading notification '%s'", notif.Name)
	if notif.If != "" {
		notif.internal.condition, err = lang.NewCondition(notif.If)
		if err != nil {
			log.Fatalf("while parsing `%s`: %s", notif.If, err)
		}
	}
}
