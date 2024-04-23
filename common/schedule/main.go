package schedule

import (
	"time"
)

type Schedule struct {
	WeeklySchedule *WeeklySchedule `yaml:"weekly_schedule" json:"weeklySchedule"`
	DailySchedule  *DailySchedule  `yaml:"daily_schedule" json:"dailySchedule"`
}

func (s *Schedule) Resolve() (Interface, error) {
	if s.WeeklySchedule != nil {
		return s.WeeklySchedule.Resolve()
	}
	if s.DailySchedule != nil {
		return s.DailySchedule.Resolve()
	}
	return nil, nil
}

type Interface interface {
	Match(*time.Time) bool
}
