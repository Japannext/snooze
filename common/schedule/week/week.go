package weekday

import (
  "time"
  "fmt"
)

var weekdays = []string{
  "Sunday", // 0
  "Monday", // 1
  "Tuesday", // 2
  "Wednesday", // 3
  "Thursday", // 4
  "Friday", // 5
  "Saturday", // 6
}

type Weekday time.Weekday

// Transform a string to a Weekday (if possible)
func Parse(w string) (Weekday, error) {
  for i, s := range weekdays {
    if s == w {
      return i, nil
    }
  }
  return 0, fmt.Errorf("'%s' is not a valid weekday", w)
}

func (w *Weekday) String() string {
  return w.String()
}

type WeekdayTimeConfig struct {
  Weekday string `yaml:"weekday" json:"weekday"`
  Time string `yaml:"time" json:"time"`
  TimeZone string `yaml:"timezone" json:"timezone"`
}

type WeeklyConfig struct {
  From WeekdayTime `yaml:"from" json:"from"`
  To WeekdayTime `yaml:"to" json:"to"`
}

type WeekdayTime struct {
  weekday time.Weekday
  hour int
  minute int
}

// Transform a weekday string to a number
func weekdayToNumber(w string) (int, error) {
  for i, s := range weekdays {
    if s == w {
      return i, nil
    }
  }
  return 0, fmt.Errorf("'%s' is not a valid weekday", w)
}

func NewWeekdayTime(cfg *WeekdayTimeConfig) (*WeekdayTime, error) {

  w : =

  return &WeekdayTime{
    weekday: weekdays[]
  }
}

type Weekly struct {
  From WeekdayTime
  To WeekdayTime
}

