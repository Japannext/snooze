package grouping

import (
  api "github.com/japannext/snooze/common/api/v2"
  "github.com/japannext/snooze/common/utils"
)

func Process(alert *api.Alert) error {
  var m map[string]string

  for _, rule := range computedRules {
    if rule.Condition.Match(alert) {
      fields := rule.GroupBy
      for _, fi := range fields {
        v, found := fi.Get(alert)
        m[fi.String()] = v
      }
      break
    }
  }

  alert.GroupHash = ComputeHash(m)
  return nil
}
