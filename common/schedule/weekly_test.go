package schedule

import (
  "time"

  "testing"
  "github.com/stretchr/testify/assert"
)

type weeklyTest struct {
  Name string
  FromWeekday string
  FromTime string
  ToWeekday string
  ToTime string
  Time string
  Expected bool
}

// 20 cases (going through all code cases in order)
// Friday     Saturday   Sunday     Monday     Tuesday    Wednesday
// 2024-04-19 2024-04-20 2024-04-21 2024-04-22 2024-04-23 2024-04-24
var weeklyTests = []weeklyTest{
  // Simple case
  {"Simple case #1: (a.hour < w.hour) && (w.hour < b.hour)",
    "Friday", "08:00", "Friday", "17:30", "2024-04-19T10:30:00+09:00", true},
  // (a.hour == w.hour) && (a.minute < w.minute)
  {"Simple case #2: (a.hour == w.hour) && (a.minute < w.minute)",
    "Friday", "08:00", "Friday", "17:30", "2024-04-19T08:01:00+09:00", true},
  {"Simple case #3: (w.hour == b.hour) && (w.minute < b.minute)",
    "Friday", "08:00", "Friday", "17:30", "2024-04-19T17:01:00+09:00", true},
  // else
  {"Simple case #4: else",
    "Friday", "08:00", "Friday", "17:30", "2024-04-19T18:01:00+09:00", false},

  // Reverse case a > b
  {"Reverse case #1: (a.weekday < w.weekday)",
    "Friday", "17:30", "Monday", "08:30", "2024-04-20T07:00:00+09:00", true},
  {"Reverse case #2: (w.weekday < b.weekday)",
    "Friday", "17:30", "Monday", "08:30", "2024-04-21T18:30:00+09:00", true},
  {"Reverse case #3: (a.hour < w.hour)",
    "Friday", "17:30", "Monday", "08:30", "2024-04-19T20:00:00+09:00", true},
  {"Reverse case #4: (w.hour == a.hour) && (a.minute < w.minute)",
    "Friday", "17:30", "Monday", "08:30", "2024-04-19T17:45:00+09:00", true},
  {"Reverse case #5: else",
    "Friday", "17:30", "Monday", "08:30", "2024-04-19T17:00:00+09:00", false},
  {"Reverse case #6: else ",
    "Friday", "17:30", "Monday", "08:30", "2024-04-22T18:00:00+09:00", false},

  // Natural order a < b
  {"Natural order #1: (a.weekday < w.weekday) && (w.weekday < b.weekday)",
    "Monday", "22:00", "Wednesday", "01:30", "2024-04-23T10:00:00+09:00", true},
  {"Natural order #2: (a.hour < w.hour)",
    "Monday", "22:00", "Wednesday", "01:30", "2024-04-22T23:59:59+09:00", true},
  {"Natural order #3: (a.hour == w.hour) && (a.minute < w.minute)",
    "Monday", "22:00", "Wednesday", "01:30", "2024-04-22T22:10:00+09:00", true},
  {"Natural order #4: (w.hour < b.hour)",
    "Monday", "22:00", "Wednesday", "01:30", "2024-04-24T00:10:00+09:00", true},
  {"Natural order #4: (w.hour == b.hour) && (w.minute <= b.minute)",
    "Monday", "22:00", "Wednesday", "01:30", "2024-04-24T01:10:00+09:00", true},
  {"Natural order #5: else",
    "Monday", "22:00", "Wednesday", "01:30", "2024-04-24T01:40:00+09:00", false},
}

func TestWeeklyMatch(t *testing.T) {

  tzName := "Asia/Tokyo"

  for _, td := range weeklyTests {
    t.Run(td.Name, func(t *testing.T) {
      wr := &WeeklyRepr{
        From: &WeekTimeRepr{td.FromWeekday, td.FromTime},
        To: &WeekTimeRepr{td.ToWeekday, td.ToTime},
        TimeZone: tzName,
      }
      w, err := wr.Resolve()
      assert.NoError(t, err)

      tt, err := time.Parse(time.RFC3339, td.Time)
      assert.NoError(t, err)

      assert.Equal(t, td.Expected, w.Match(&tt))
    })
  }
}
