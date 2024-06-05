package schedule

import (
	"time"
)

type DailySchedule struct {
	From     string `yaml:"from" json:"from"`
	To       string `yaml:"to" json:"to"`
	TimeZone string `yaml:"timezone" json:"timezone"`
}

type hourminute struct {
	hour   int
	minute int
}

type Daily struct {
	From     *hourminute
	To       *hourminute
	TimeZone *time.Location
}

func (dr *DailySchedule) Resolve() (*Daily, error) {
	h1, m1, err := parseTime(dr.From)
	if err != nil {
		return nil, err
	}
	h2, m2, err := parseTime(dr.To)
	if err != nil {
		return nil, err
	}
	tz, err := time.LoadLocation(dr.TimeZone)
	if err != nil {
		return nil, err
	}
	return &Daily{
		From:     &hourminute{h1, m1},
		To:       &hourminute{h2, m2},
		TimeZone: tz,
	}, nil
}

func (hm *hourminute) before(a *hourminute) bool {
	if hm.hour < a.hour {
		return true
	}
	if hm.hour == a.hour && hm.minute < a.minute {
		return true
	}
	return false
}

func (hm *hourminute) after(b *hourminute) bool {
	return b.before(hm)
}

func (hm *hourminute) InBetween(a, b *hourminute) bool {
	// Reverse case
	if a.before(b) {
		if hm.after(a) || hm.before(b) {
			return true
		}
		return false
	}

	// Natural order
	if hm.after(a) && hm.before(b) {
		return true
	}

	return false
}

func (s *Daily) Match(t *time.Time) bool {
	tt := t.In(s.TimeZone)
	hm := &hourminute{tt.Hour(), tt.Minute()}
	return hm.InBetween(s.From, s.To)
}
