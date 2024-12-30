package schedule

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Schedule struct {
	Always bool            `json:"always"           yaml:"always"`
	Weekly *WeeklySchedule `json:"weekly,omitempty" yaml:"weekly"`
	Daily  *DailySchedule  `json:"daily,omitempty"  yaml:"daily"`

	internal struct {
		schedule Interface
	}
}

func (s *Schedule) Load() {
	switch {
	case s.Always:
		s.internal.schedule = &AlwaysSchedule{}
	case s.Weekly != nil:
		s.internal.schedule = s.Weekly
	case s.Daily != nil:
		s.internal.schedule = s.Daily
	default:
		log.Fatalf("empty schedule defined")
	}

	s.internal.schedule.Load()
}

func Default() *Schedule {
	return &Schedule{Always: true}
}

func (s *Schedule) Match(t *time.Time) bool {
	return s.internal.schedule.Match(t)
}

type Interface interface {
	Load()
	Match(*time.Time) bool
}
