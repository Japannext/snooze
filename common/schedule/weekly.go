package schedule

import (
  "fmt"
  "time"
)

// Like golang's time.Weekday
var weekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// Transform a string to a Weekday (if possible)
func parseWeekday(w string) (int, error) {
  for i, s := range weekdays {
    if s == w {
      return i, nil
    }
  }
  return 0, fmt.Errorf("'%s' is not a valid weekday", w)
}

// Parse a time string in the format "HH:MM" to an hour and minute integers
func parseTime(s string) (int, int, error) {
  var h, m int
  t, err := time.Parse("15:04", s)
  if err !=  nil {
    return h, m, err
  }
  return t.Hour(), t.Minute(), nil
}

type WeeklySchedule struct {
  From *WeekTimeRepr `yaml:"from" json:"from"`
  To *WeekTimeRepr `yaml:"to" json:"to"`
  TimeZone string `yaml:"timezone" json:"timezone"`
}

type WeekTimeRepr struct {
  Weekday string `yaml:"weekday" json:"weekday"`
  Time string `yaml:"time" json:"time"`
}

type weektime struct {
  weekday int
  hour int
  minute int
}

func compareInt(a, b int) int {
  if a == b {
    return 0
  }
  if a < b {
    return -1
  }
  return +1
}

func isWeekdayInBetween(x, a, b int) bool {
  if (a > b) {
    return (x <= b) || (a <= x)
  }
  return (a <= x) && (x <= b)
}

// Check if a weektime is between 2 weektimes
func (w *weektime) InBetween(a, b *weektime) bool  {
  // Simple case
  if (a.weekday == b.weekday) {
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
  if (a.weekday > b.weekday) {
    if (w.weekday < b.weekday) || (a.weekday < w.weekday) {
      return true
    }
    if (a.weekday == w.weekday) {
      if (a.hour < w.hour) {
        return true
      }
      if (w.hour == a.hour) && (a.minute < w.minute) {
        return true
      }
      return false
    }
    if (w.weekday == b.weekday) {
      if (w.hour < b.hour) {
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
  if (a.weekday == w.weekday) {
    if (a.hour < w.hour) {
      return true
    }
    if (a.hour == w.hour) && (a.minute <= w.minute) {
      return true
    }
    return false
  }
  if (w.weekday == b.weekday) {
    if (w.hour < b.hour) {
      return true
    }
    if (w.hour == b.hour) && (w.minute <= b.minute) {
      return true
    }
    return false
  }

  return false
}

func (wtr *WeekTimeRepr) toWeekTime() (*weektime, error) {
  hour, minute, err := parseTime(wtr.Time)
  if err != nil {
    return nil, err
  }
  weekday, err := parseWeekday(wtr.Weekday)
  if err != nil {
    return nil, err
  }
  return &weektime{weekday, hour, minute}, nil
}

type Weekly struct {
  from *weektime
  to *weektime
  tz *time.Location
}

func (w *WeeklySchedule) Resolve() (*Weekly, error) {
  from, err := w.From.toWeekTime()
  if err != nil {
    return nil, err
  }
  to, err := w.To.toWeekTime()
  if err != nil {
    return nil, err
  }

  // Edge-cases
  if from.weekday == to.weekday {
    if to.hour < from.hour {
      return nil, fmt.Errorf("same day, but hours backwards!")
    }
    if from.hour == to.hour && to.minute < from.minute {
      return nil, fmt.Errorf("same day, same hour, but minutes backwards!")
    }
  }

  tz, err := time.LoadLocation(w.TimeZone)
  if err != nil {
    return nil, err
  }

  return &Weekly{from, to, tz}, nil
}

func (s *Weekly) Match(t *time.Time) bool {
  tt := t.In(s.tz)
  w := &weektime{int(tt.Weekday()), tt.Hour(), tt.Minute()}
  return w.InBetween(s.from, s.to)
}
