export type Schedule = {
  always?: boolean
  weekly?: WeeklySchedule
  daily?: DailySchedule
}

export type WeeklySchedule = {
  from: Weektime
  to: Weektime
  timezone: string
}

export type Weektime = {
  // weekday by name (e.g. "Monday", "Tuesday", ...)
  weekday: string
  // Time in a hour:minute format
  time: string
}

export type DailySchedule = {
  from: string
  to: string
  timezone: string
}
