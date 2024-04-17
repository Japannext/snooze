package notification

import (
  set "github.com/deckarep/golang-set/v2"

  "github.com/japannext/snooze/common/utils"
)

func Process(alert *v2.Alert) error {

  if alert.Internals.Snoozed {
    return nil
  }

  var queues = set.NewSet[*NotificationQueue]()
  var merr = utils.NewMultiError("Failed to notify alert trace_id=%s.")

  for _, rule := range computedRules {
    if rule.Condition.Test(alert) {
      for _, q := range rule.Queues {
        queues.Add(q)
      }
      if rule.StopThere {
        break
      }
    }
  }

  for _, q := range queues.Iter() {
    notif := api.Notification{}
    if err := q.Publish(ctx, notif); err != nil {
      merr.AppendErr(err)
      continue
    }
  }

  if merr.HasErrors() {
    return merr
  }

  return nil
}
