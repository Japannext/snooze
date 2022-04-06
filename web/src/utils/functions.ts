/**
 * A list of functions without big dependencies on Vue
**/

import moment from 'moment'

/** Convert an ISO timestamp into a relative time to display to the user
 * @param date {string} The date in ISO format (timestamp from database objects)
 * @param showSeconds {boolean} Show seconds in the trimmed date
 * @returns {string} A string displaying a relative time, in pretty print format
**/
export function prettyDate(date: string, showSeconds: boolean): string {
  if (!date) {
    return 'Empty'
  }
  let mDate = moment(date)
  let newDate = ''
  const now = moment()
  let hours_only = false
  if (!mDate.isValid()) {
    mDate = moment(date, ['HH:mm'])
    hours_only = true
  }
  if (mDate.year() == now.year()) {
    if (mDate.format('MM-DD') == now.format('MM-DD')) {
      if (hours_only) {
        newDate = mDate.format('HH:mm')
      } else if (showSeconds) {
        newDate = 'Today' + ' ' + mDate.format('HH:mm:ss')
      } else {
        newDate = 'Today' + ' ' + mDate.format('HH:mm')
      }
    } else {
      if (showSeconds) {
        newDate = mDate.format('MMM Do HH:mm:ss')
      } else {
        newDate = mDate.format('MMM Do HH:mm')
      }
    }
  } else {
    newDate = mDate.format('MMM Do YYYY')
  }
  return newDate
}

/** Transform durations in seconds to a human readable value
 * @param seconds Number of seconds of the duration to print
 * @returns A human readable display of the duration
**/
export function prettyDuration(seconds: string | number): string {
  let secondsNumber: number
  if (typeof seconds == 'number') {
    secondsNumber = seconds
  } else {
    secondsNumber = parseInt(seconds, 10)
  }
  if (secondsNumber < 0) {
    return '0s'
  }

  const days    = Math.floor(secondsNumber / (3600*24))
  const hours   = Math.floor(secondsNumber / 3600) % 24
  const minutes = Math.floor(secondsNumber / 60) % 60
  const secs    = secondsNumber % 60

  let output = ''
  if (days > 0) {
    output += days + 'd '
  }
  if (hours > 0) {
    output += hours + 'h '
  }
  if (minutes > 0) {
    output += minutes + 'm '
  }
  if (secs > 0) {
    output += secs + 's'
  }
  return output
}

/** Pretty display for numbers
 * @param x Number to pretty print
 * @returns A more human readable number (a comma every 3 digits)
**/
export function prettyNumber(x: number): string {
  if (isNaN(x)) {
    return '0'
  }
  return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

/** A pretty print for a time to live
 * @param seconds The number of seconds left for a resource to live
 * @returns A time to live in hh::mm::ss format
**/
export function prettyTTL(seconds: string): string {
  var secondsNumber = parseInt(seconds, 10)
  var hours   = Math.floor(secondsNumber / 3600)
  var minutes = Math.floor(secondsNumber / 60) % 60
  var secs = secondsNumber % 60

  const hoursString = String(hours).padStart(2, '0')
  const minutesString = String(minutes).padStart(2, '0')
  const secsString = String(secs).padStart(2, '0')
  if (secs < 0) {
    return '00:00:00'
  } else {
    return hoursString + ':' + minutesString + ':' + secsString
  }
}

/** Capitalize the first letter of a given string
 * @param string A given string to capitalize
 * @returns A capitalized string
**/
export function capitalizeFirstLetter(string: string): string {
  return string.charAt(0).toUpperCase() + string.slice(1)
}


export function truncateMessage(message: string, size: number = 280): string {
  if (message == null || message.length <= size) {
    return message
  }
  return message.slice(0, size) + '...'
}

export function more(data, limit=5) {
  var lines = data.split(/\r?\n/)
  if (lines.length > limit) {
    var start = lines.slice(0, limit).join('\n')
    var end = lines.slice(limit+1, -1).join('\n')
    return [start, end]
  } else {
    return [data, '']
  }
}

