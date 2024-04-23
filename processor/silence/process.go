package silence

import (
  api "github.com/japannext/snooze/common/api/v2"
)

func Process(alert *api.Alert) error {

  for _, rule := range computedRules {
    if rule.Condition.Test(alert) {
      // Silence the alert
      alert.Mute.Enabled = true
      alert.Mute.Component = "silence"
      alert.Mute.Rule = rule.String()
      alert.Mute.SkipNotification = true
      break
    }
  }

  return nil
}
