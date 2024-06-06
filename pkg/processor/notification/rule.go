package notification

import (
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
)

type Rule struct {
	If       string   `yaml:"if"`
	Channels []string `yaml:"channels"`
}

type computedRule struct {
	Condition *lang.Condition
	Queues    []*rabbitmq.NotificationQueue
}

func compute(rule *Rule) *computedRule {
	var condition *lang.Condition
	var err error
	if rule.If != "" {
		condition, err = lang.NewCondition(rule.If)
		if err != nil {
			log.Fatalf("while parsing `%s`: %s", rule.If, err)
		}
	}

	var qq []*rabbitmq.NotificationQueue

	for _, name := range rule.Channels {
		q, found := queues[name]
		if !found {
			q = ch.NewQueue(name)
			queues[name] = q
		}
		qq = append(qq, q)
	}
	return &computedRule{
		Condition: condition,
		Queues:    qq,
	}
}

var ch *rabbitmq.NotificationChannel
var queues = make(map[string]*rabbitmq.NotificationQueue)

var computedRules []*computedRule

var log *logrus.Entry

func InitRules(rules []*Rule, defaults []string) {
	log = logrus.WithFields(logrus.Fields{"module": "notification"})
	ch = rabbitmq.InitNotificationChannel()

	for _, rule := range rules {
		computedRules = append(computedRules, compute(rule))
	}
	computedRules = append(computedRules, compute(&Rule{Channels: defaults}))
}
