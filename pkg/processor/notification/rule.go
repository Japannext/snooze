package notification

import (
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/utils"
)

type Rule struct {
	If       string   `yaml:"if"`
	Channels []string `yaml:"channels"`

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
var log *logrus.Entry
var producers map[string]*rabbitmq.Producer

func Startup(rules []*Rule, defaults []string) {
	log = logrus.WithFields(logrus.Fields{"module": "notification"})

	rules = append(rules, &Rule{Channels: defaults})

	var queues = utils.NewOrderedStringSet()

	for _, rule := range rules {
		rule.Startup()
		computedRules = append(computedRules, rule)
		queues.AppendMany(rule.Channels)
	}

	for _, queue := range queues.Items() {
		producers[queue] = rabbitmq.NewNotificationProducer(queue)
	}
}
