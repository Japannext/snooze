package schedule

import (
	"time"
	"fmt"
)

type Schedule struct {
	Always bool            `json:"always"           yaml:"always"`
	Weekly *WeeklySchedule `json:"weekly,omitempty" yaml:"weekly"`
	Daily  *DailySchedule  `json:"daily,omitempty"  yaml:"daily"`

	internal struct {
		schedule Interface
	}
}

func (s *Schedule) Load() error {
	switch {
	case s.Always:
		s.internal.schedule = &AlwaysSchedule{}
	case s.Weekly != nil:
		s.internal.schedule = s.Weekly
	case s.Daily != nil:
		s.internal.schedule = s.Daily
	default:
		return fmt.Errorf("empty schedule")
	}

	return s.internal.schedule.Load()
}

func Default() *Schedule {
	return &Schedule{Always: true}
}

func (s *Schedule) Match(t *time.Time) bool {
	return s.internal.schedule.Match(t)
}

type Interface interface {
	Load() error
	Match(*time.Time) bool
}
