package schedule

import (
  "time"
  "strconv"
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

type WeeklyRepr struct {
  From WeekTimeRepr `yaml:"from" json:"from"`
  To WeekTimeRepr `yaml:"to" json:"to"`
}

type WeekTimeRepr struct {
  Weekday string `yaml:"weekday" json:"weekday"`
  Time string `yaml:"time" json:"time"`
  TimeZone string `yaml:"timezone" json:"timezone"`
}

type WeekTime struct {
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

// Compare 2 weekdays
func (w1 *WeekTime) Compare(w2 *WeekTime) int {
  w := compareInt(w1.weekday, w2.weekday)
  if w == 0 {
    h := compareInt(w1.hour, w2.hour)
    if h == 0 {
      m := compareInt(w1.minute, w2.minute)
      return m
    }
    return h
  }
  return w
}

func (w1 *WeekTime) Before(w2 *WeekTime) bool {
  
}

func getWeekTime(cfg *WeekTimeRepr) (*WeekTime, error) {
  hour, minute, err := parseTime(cfg.Time)
  if err != nil {
    return nil, err
  }
  weekday, err := parseWeekday(cfg.Weekday)
  if err != nil {
    return nil, err
  }
  return &WeekTime{weekday, hour, minute}, nil
}

type Weekly struct {
  From WeekTime
  To WeekTime
}

func getWeekly(w *WeeklyRepr) (*Weekly, error) {
  from, err := getWeekTime(w.From)
  if err != nil {
    return nil, err
  }
  to := getWeekTime(w.To)
  if err != nil {
    return nil, err
  }
  return &Weekly{from, to}, nil
}

func (s *Weekly) MatchNow() bool {
  return s.Match(time.Now())
}

func (s *Weekly) Match(t *time.Time) bool {
  wt := &WeekTime{t.Weekday(), t.Hour(), t.Minute()}

  if wt.Compare(s.From) >= 0 {
  }

  if w < s.From.Weekday {
    return false
  }
  if w > s.To.Weekday {
    return false
  }
  if w == s.From.Weekday && 
}
