package notification

import (
  "github.com/japannext/snooze/common/condition"
)

type Rule struct {
  If string `yaml:"if"`
  Channels []string `yaml:"channels"`
}

type computedRule struct {
  Condition *condition.Condition
  Queues []*rabbitmq.NotificationQueue
}

var computedRules []computedRule

func compute(rule *Rule) *computedRule {
  c, err := condition.Parse(rule.If)
  if err != nil {
    log.Fatal(err)
  }
  var queues []*rabbitmq.NotificationQueue
  for _, ch := range rule.Channels {
    q := InitNotificationQueue(ch)
    queues = append(queues, q)
  }
  return &computedRule{
    Condition: c,
    Queues: queues,
  }
}

func InitRules(rules []Rule, defaults []string) {
  for _, rule := range rules {
    computedRules = append(computedRules, compute(rule))
  }
  computedRules = append(computedRules, &computedRule{Channels: defaults})
}
