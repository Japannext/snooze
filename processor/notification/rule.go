package notification

import (
	"github.com/sirupsen/logrus"

	"github.com/PaesslerAG/gval"
	"github.com/japannext/snooze/common/rabbitmq"
)

type Rule struct {
	If       string   `yaml:"if"`
	Channels []string `yaml:"channels"`
}

type computedRule struct {
	Matcher gval.Evaluable
	Queues  []*rabbitmq.NotificationQueue
}

func compute(rule *Rule) *computedRule {
	matcher, err := gval.Full().NewEvaluable(rule.If)
	if err != nil {
		log.Fatal(err)
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
		Matcher: matcher,
		Queues:  qq,
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
