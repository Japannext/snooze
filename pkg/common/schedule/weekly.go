package schedule

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Like golang's time.Weekday.
var weekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// Transform a string to a Weekday (if possible).
func parseWeekday(w string) (int, error) {
	for i, s := range weekdays {
		if s == w {
			return i, nil
		}
	}
	return 0, fmt.Errorf("'%s' is not a valid weekday", w)
}

// Parse a time string in the format "HH:MM" to an hour and minute integers.
func parseTime(s string) (int, int, error) {
	var h, m int
	t, err := time.Parse("15:04", s)
	if err != nil {
		return h, m, err
	}
	return t.Hour(), t.Minute(), nil
}

type WeeklySchedule struct {
	From     *Weektime `json:"from"     yaml:"from"`
	To       *Weektime `json:"to"       yaml:"to"`
	TimeZone string    `json:"timezone" yaml:"timezone"`

	internal struct {
		timezone *time.Location
	}
}

func (s *WeeklySchedule) Load() {
	s.From.Load()
	s.To.Load()
	var err error
	s.internal.timezone, err = time.LoadLocation(s.TimeZone)
	if err != nil {
		log.Fatalf("failed to load timezone `%s`: %s", s.TimeZone, err)
	}

	from := s.From.internal
	to := s.To.internal
	// Handling edge-case
	if from.weekday == to.weekday {
		if to.hour < from.hour {
			log.Fatalf("error defining weekly schedule: same day, but hours backwards!")
		}
		if to.hour == from.hour && to.minute < from.minute {
			log.Fatalf("error defining weekly schedule: same day, same hour, but minutes backwards!")
		}
	}
}

func (s *WeeklySchedule) Match(t *time.Time) bool {
	tt := t.In(s.internal.timezone)
	wt := NewWeektimeFromTime(&tt)
	return wt.InBetween(s.From, s.To)
}

// A specific time within a week, represented by a weekday, and a time (hour:minute).
type Weektime struct {
	Weekday string `json:"weekday" yaml:"weekday"`
	Time    string `json:"time"    yaml:"time"`

	internal struct {
		weekday int
		hour    int
		minute  int
	}
}

func NewWeektimeFromTime(t *time.Time) *Weektime {
	wt := &Weektime{}
	wt.internal.weekday = int(t.Weekday())
	wt.internal.hour = t.Hour()
	wt.internal.minute = t.Minute()
	return wt
}

func (wt *Weektime) Load() {
	var err error
	wt.internal.hour, wt.internal.minute, err = parseTime(wt.Time)
	if err != nil {
		log.Fatalf("failed to parse time `%s`", wt.Time)
	}
	wt.internal.weekday, err = parseWeekday(wt.Weekday)
	if err != nil {
		log.Fatalf("failed to parse weekday `%s`", wt.Weekday)
	}
}

// Check if a weektime is between 2 weektimes.
func (wt *Weektime) InBetween(wa, wb *Weektime) bool {
	a := wa.internal
	b := wb.internal
	w := wt.internal
	// Simple case
	if a.weekday == b.weekday {
		if (a.hour < w.hour) && (w.hour < b.hour) {
			return true
		}
		if (a.hour == w.hour) && (a.minute < w.minute) {
			return true
		}
		if (w.hour == b.hour) && (w.minute < b.minute) {
			return true
		}
		return false
	}

	// Reverse case a > b
	if a.weekday > b.weekday {
		if (w.weekday < b.weekday) || (a.weekday < w.weekday) {
			return true
		}
		if a.weekday == w.weekday {
			if a.hour < w.hour {
				return true
			}
			if (w.hour == a.hour) && (a.minute < w.minute) {
				return true
			}
			return false
		}
		if w.weekday == b.weekday {
			if w.hour < b.hour {
				return true
			}
			if (w.hour == b.hour) && (w.minute < b.minute) {
				return true
			}
			return false
		}
	}

	// Natural order a < b
	if (a.weekday < w.weekday) && (w.weekday < b.weekday) {
		return true
	}
	if a.weekday == w.weekday {
		if a.hour < w.hour {
			return true
		}
		if (a.hour == w.hour) && (a.minute <= w.minute) {
			return true
		}
		return false
	}
	if w.weekday == b.weekday {
		if w.hour < b.hour {
			return true
		}
		if (w.hour == b.hour) && (w.minute <= b.minute) {
			return true
		}
		return false
	}

	return false
}
