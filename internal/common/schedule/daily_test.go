package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dailyTest struct {
	Name     string
	FromTime string
	ToTime   string
	Time     string
	Expected bool
}

var dailyTests = []dailyTest{
	{
		"Natural order #1",
		"08:00", "17:00", "2024-04-19T10:30:00+09:00", true,
	},
	{
		"Natural order #2",
		"08:00", "17:00", "2024-04-19T07:30:00+09:00", true,
	},
	{
		"Natural order #3",
		"08:00", "17:00", "2024-04-19T07:30:00+09:00", true,
	},
}

func TestMatch(t *testing.T) {
	tzName := "Asia/Tokyo"

	for _, td := range dailyTests {
		t.Run(td.Name, func(t *testing.T) {
			s := &DailySchedule{}
			s.From = td.FromTime
			s.To = td.ToTime
			s.TimeZone = tzName
			s.Load()

			tt, err := time.Parse(time.RFC3339, td.Time)
			assert.NoError(t, err)

			assert.Equal(t, td.Expected, s.Match(&tt))
		})
	}
}
