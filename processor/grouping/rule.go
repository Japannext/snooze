package grouping

import (
  log "github.com/sirupsen/logrus"
  "github.com/japannext/snooze/common/condition"
  "github.com/japannext/snooze/common/field"
)

type Rule struct {
  If string `yaml:"if"`
  GroupBy []string `yaml:"group_by"`
}

type computedRule struct {
  Condition condition.Interface
  Fields []*field.AlertField
}

var computedRules []*computedRule

func compute(rule *Rule) *computedRule {
  c, err := condition.Parse(rule.If)
  if err != nil {
    log.Fatal(err)
  }
  var fields []*field.AlertField
  for _, f := range rule.GroupBy {
    fi, err := field.Parse(f)
    if err != nil {
      log.Fatal(err)
    }
    fields = append(fields, fi)
  }
  return &computedRule{
    Condition: c.Resolve(),
    Fields: fields,
  }
}

func InitRules(rules []*Rule) {
  for _, rule := range rules {
    computedRules = append(computedRules, compute(rule))
  }
}
