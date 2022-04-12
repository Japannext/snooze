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


/** Truncate a message after a given number of characters, and suffix it with '...'
 * @param message The string to truncate
 * @param size The size limit in number of characters
**/
export function truncateMessage(message: string|null, size=280): string|null {
  if (message == null || message.length <= size) {
    return message
  }
  return message.slice(0, size) + '...'
}

/** Separate a message into a truncated version, and the rest, based of a number of lines.
 * @param data The string to truncate
 * @param limit The number of lines after which we will truncate
 * @returns The lines before truncate, and after the truncate
**/
export function more(data: string, limit=5): [string, string] {
  const lines = data.split(/\r?\n/)
  if (lines.length > limit) {
    const start = lines.slice(0, limit).join('\n')
    const end = lines.slice(limit+1, -1).join('\n')
    return [start, end]
  } else {
    return [data, '']
  }
}

/** Parse an input, and try to determine if it's a string or a number
 * @param data The value to parse
 * @returns The parsed value
**/
export function inputParser(data: string): number|string {

  const dataAsNumber = Number(data)

  // Handling the quote case (dequote it)
  if (/^".*"$/.test(data) || /^'.*'$/.test(data)) {
    const subString = data.substr(1, data.length -2)
    return subString
  }

  if (isNaN(dataAsNumber)) {
    return data
  } else {
    return dataAsNumber
  }
}

/** Get a number or string from the database, and return a string compatible
 * with inputParse.
 * @param data The data that should result of inputParser
 * @returns A string compatible with inputParser
**/
export function inputDeParser(data: number|string): string {
  switch(typeof data) {
    case 'number': {
      return String(data)
    }
    case 'string': {
      const dataAsNumber = Number(data)
      if (data == '' || isNaN(dataAsNumber)) {
        return data
      } else {
        return `"${dataAsNumber}"`
      }
    }
  }
}

/** Return the human readable name for weekday by number
 * @param index The index of weekday (starting at 0 from Sunday)
 * @returns The corresponding human readable weekday
**/
export function getWeekday(index: number): string {
  switch(index) {
    case 0:
      return 'Sunday'
    case 1:
      return 'Monday'
    case 2:
      return 'Tuesday'
    case 3:
      return 'Wednesday'
    case 4:
      return 'Thursday'
    case 5:
      return 'Friday'
    case 6:
      return 'Saturday'
    default:
      return 'Invalid weekday ' + String(index)
  }
}

/** Stop an event, and its propagation
 * @param event The event to stop
 * @param {object} obj
 * @param {boolean} obj.preventDefault Prevent the default propagation of an event
 * @param {boolean} obj.propagation Stop the propagation of an event
 * @param {boolean} obj.immediatePropagation Stop the immediate propagation of an event
**/
export function stopEvent(event: Event, {preventDefault=true, propagation=true, immediatePropagation=false} = {}) {
  if (preventDefault) {
    event.preventDefault()
  }
  if (propagation) {
    event.stopPropagation()
  }
  if (immediatePropagation) {
    event.stopImmediatePropagation()
  }
}

/** Copy a given string to the user's clipboard
 * @param {string} text The text to copy to clipboard
**/
export function copyToClipBoard(text: string) {
  navigator.clipboard.writeText(text)
  .catch(err => {
    console.error('Could not copy text to clipboard: ', err)
  })
}

