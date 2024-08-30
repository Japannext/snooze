package schedule

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type DailySchedule struct {
	From     string `yaml:"from" json:"from"`
	To       string `yaml:"to" json:"to"`
	TimeZone string `yaml:"timezone" json:"timezone"`

	internal struct {
		from *hourminute
		to *hourminute
		timezone *time.Location
	}
}

func (dr *DailySchedule) Load() {
	h1, m1, err := parseTime(dr.From)
	if err != nil {
		log.Fatalf("failed to load time `%s`: %s", dr.From, err)
	}
	dr.internal.from = &hourminute{h1, m1}
	h2, m2, err := parseTime(dr.To)
	if err != nil {
		log.Fatalf("failed to load time `%s`: %s", dr.To, err)
	}
	dr.internal.to = &hourminute{h2, m2}
	tz, err := time.LoadLocation(dr.TimeZone)
	if err != nil {
		log.Fatalf("failed to load timezone `%s`: %s", dr.TimeZone, err)
	}
	dr.internal.timezone = tz
}

type hourminute struct {
	hour   int
	minute int
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

func (s *DailySchedule) Match(t *time.Time) bool {
	tt := t.In(s.internal.timezone)
	hm := &hourminute{tt.Hour(), tt.Minute()}
	return hm.InBetween(s.internal.from, s.internal.to)
}
