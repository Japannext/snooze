package silence

import (
  log "github.com/sirupsen/logrus"
  "github.com/japannext/snooze/common/condition"
  "github.com/japannext/snooze/common/schedule"
)

type Rule struct {
  Name string `yaml:"name"`
  If string `yaml:"if"`
  Schedule *schedule.Schedule `yaml:",inline"`
}

type computedRule struct {
  Name string
  Condition condition.Interface
  Schedule schedule.Interface
}

func (r *computedRule) String() string {
  if r.Name != "" {
    return r.Name
  }
  return r.Condition.String()
}

var computedRules []*computedRule

func compute(rule *Rule) *computedRule {
  c, err := condition.Parse(rule.If)
  if err != nil {
    log.Fatal(err)
  }
  s, err := rule.Schedule.Resolve()
  if err != nil {
    log.Fatal(err)
  }
  return &computedRule{
    Name: rule.Name,
    Condition: c.Resolve(),
    Schedule: s,
  }
}

func InitRules(rules []*Rule) {
  for _, rule := range rules {
    computedRules = append(computedRules, compute(rule))
  }
}
