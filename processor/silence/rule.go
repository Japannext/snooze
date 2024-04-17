package grouping

import (
  "crypto/md5"
  "sort"

  api "github.com/japannext/snooze/common/api/v2"
  "github.com/japannext/snooze/common/condition"
  "github.com/japannext/snooze/common/schedule"
  "github.com/japannext/snooze/common/field"
  "github.com/japannext/snooze/common/utils"
)

type Rule struct {
  If string `yaml:"if"`
  WeeklySchedule *schedule.Weekly `yaml:"weekly_schedule"`
  DailySchedule *schedule.Daily `yaml:"daily_schedule"`
}

type computedRule struct {
  Condition *condition.Condition
  Fields []field.AlertField
}

var computedRules []ComputedRule

func compute(r *Rule) *computedRule {
  c, err := condition.Parse(rule.If)
  if err != nil {
    log.Fatal(err)
  }
  var fields []field.AlertField
  for _, f := range rule.GroupBy {
    fi, err := field.Parse(f)
    if err != nil {
      log.Fatal(err)
    }
  }
  return &ComputedRule{
    Condition: c,
    Fields: fields,
  }
}

func InitRules(rules []Rule) {
  for _, rule := range rules {
    computedRules = append(computedRules, compute(rule))
  }
}
