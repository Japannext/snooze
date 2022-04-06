import moment from 'moment'

export function trimDate(date: string, showSeconds: boolean): string {
  if (!date) {
    return 'Empty'
  }
  var mDate = moment(date)
  var newDate = ''
  var now = moment()
  var hours_only = false
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


