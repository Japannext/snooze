package transform

import (

)

func Process(alert *api.Alert) error {
  for _, rule := range computedRules {
    if rule.Condition.Match(alert) {
      // ...
    }
  }

  return nil
}
