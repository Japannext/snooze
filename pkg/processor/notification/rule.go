package notification

import (
	"github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/utils"
)

type Rule struct {
	If       string   `yaml:"if"`
	Destinations []api.Destination `yaml:"destinations"`

	internal struct {
		condition *lang.Condition
	}
}

func (rule *Rule) Startup() {
	var err error
	if rule.If != "" {
		rule.internal.condition, err = lang.NewCondition(rule.If)
		if err != nil {
			log.Fatalf("while parsing `%s`: %s", rule.If, err)
		}
	}
}

var computedRules []*Rule
var defaultDestinations []api.Destination
var log *logrus.Entry
var producers = map[string]*rabbitmq.Producer{}

func Startup(rules []*Rule, defaults []api.Destination) {
	log = logrus.WithFields(logrus.Fields{"module": "notification"})

	defaultDestinations = defaults

	var queues = utils.NewOrderedSet[string]()

	for _, rule := range rules {
		rule.Startup()
		computedRules = append(computedRules, rule)
		for _, dest := range rule.Destinations {
			queues.Append(dest.Queue)
		}
	}

	for _, queue := range queues.Items() {
		producers[queue] = rabbitmq.NewNotificationProducer(queue)
	}
}
