package schedule

import (
	"time"
)

type AlwaysSchedule struct{}

func (s *AlwaysSchedule) Load() error {
	return nil
}

func (s *AlwaysSchedule) Match(t *time.Time) bool {
	return true
}
