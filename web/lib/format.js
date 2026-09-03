import {
  parseISO,
  addDays,
  addMonths,
  addYears,
  isBefore,
  formatDistanceToNow,
  format,
  isValid,
} from 'date-fns'

function parseDateInput(value) {
  const date =
    value instanceof Date
      ? new Date(value.getTime())
      : typeof value === 'number'
        ? new Date(value)
        : typeof value === 'string' && value.trim()
          ? parseISO(value)
          : new Date(Number.NaN)

  return isValid(date) ? date : null
}
export function amount(n, currency = '') {
  let sign = ''
  if (currency) {
    switch (currency.toLocaleUpperCase()) {
      case 'USD':
        sign = '$'
        break
    }
  }

  return `${sign}${Number(n).toFixed(2)}`
}

export function numberWithCommas(n) {
  return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

export function dateFromNow(date) {
  const d = parseDateInput(date)
  if (!d) return INVALID_DATE

  const now = new Date()
  const oneDayLater = addDays(d, 1)
  const oneMonthLater = addMonths(oneDayLater, 1)
  const oneYearLater = addYears(oneMonthLater, 1)

  if (isBefore(now, oneDayLater)) {
    return formatDistanceToNow(d, { addSuffix: true })
  }

  if (isBefore(now, oneYearLater)) {
    return format(d, 'MMM dd')
  }

  return format(d, 'MMM dd, yyyy')
}

export function daysFromNow(d) {
  const date = parseDateInput(d)
  if (!date) return INVALID_DATE

  const nowUnix = Math.floor(Date.now() / 1000)
  const dateUnix = Math.floor(date.getTime() / 1000)
  const diffDays = ((nowUnix - dateUnix) / 86_400).toFixed()

  if (Number(diffDays) > 20 && Number(diffDays) < 60) {
    return `${diffDays} days ago`
  }

  return formatDistanceToNow(date, { addSuffix: true })
}

export function dateCalendar(date) {
  const d = parseDateInput(date)
  return d ? format(d, 'MMMM dd, yyyy') : INVALID_DATE
}

export function dateTime(date) {
  const d = parseDateInput(date)
  return d ? format(d, 'MMM dd, yyyy - h:mm a') : INVALID_DATE
}

export function dateTimeFull(date) {
  const d = parseDateInput(date)
  return d ? format(d, 'MMMM dd, yyyy - h:mm a') : INVALID_DATE
}

export function relativeFromNow(value) {
  const date = parseDateInput(value)

  return date ? formatDistanceToNow(date, { addSuffix: true }) : 'Invalid date'
}

export function errorSimple(error) {
  if (!error) {
    return ''
  }

  return error.split(':')[0]
}
