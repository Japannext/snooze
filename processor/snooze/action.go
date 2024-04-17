package snooze

import (
  "github.com/japannext/snooze/common/api/v2"
)

type SnoozeAction struct {
}

func (a *SnoozeAction) Process(item v2.Alert) error {
  return nil
}
