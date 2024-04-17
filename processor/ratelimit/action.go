package ratelimit

import (
  "github.com/japannext/snooze/common/api/v2"
)

type RateLimitAction struct {
}

func (a *RateLimitAction) Process(item v2.Alert) error {
  return nil
}
