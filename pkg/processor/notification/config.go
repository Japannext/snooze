package notification

import (
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
)

type Notification struct {
	If       string   `yaml:"if"`
	Destinations []models.Destination `yaml:"destinations"`

	internal struct {
		condition *lang.Condition
	}
}

func (notif *Notification) Load() {
	var err error
	if notif.If != "" {
		notif.internal.condition, err = lang.NewCondition(notif.If)
		if err != nil {
			log.Fatalf("while parsing `%s`: %s", notif.If, err)
		}
	}
}
