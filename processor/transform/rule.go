package transform

import (
  api "github.com/japannext/snooze/common/api/v2"
  "github.com/japannext/snooze/common/condition"
)

type Rule struct {
  If string `yaml:"if"`
}

type computeRule struct {
  Condition *condition.Condition
}

func InitRules(rules []Rule) {
  // ...
}
