package schedule

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Schedule struct {
	Always bool            `yaml:"always" json:"always"`
	Weekly *WeeklySchedule `yaml:"weekly" json:"weekly,omitempty"`
	Daily  *DailySchedule  `yaml:"daily" json:"daily,omitempty"`

	internal struct {
		schedule ScheduleInterface
	}
}

func (s *Schedule) Load() {
	switch {
	case s.Always:
		s.internal.schedule = &AlwaysSchedule{}
		break
	case s.Weekly != nil:
		s.internal.schedule = s.Weekly
		break
	case s.Daily != nil:
		s.internal.schedule = s.Daily
		break
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

type ScheduleInterface interface {
	Load()
	Match(*time.Time) bool
}
