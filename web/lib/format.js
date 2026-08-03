const MONTHS = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec',
]

const MONTHS_FULL = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
]

const MS_PER_DAY = 86400000

const toDate = value => {
  if (value == null || value === '') {
    return null
  }
  const d = new Date(value)
  return isNaN(d.getTime()) ? null : d
}

const pad2 = n => String(n).padStart(2, '0')

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

// "5 minutes ago" / "in 3 hours" style relative time (moment.fromNow equivalent).
export function fromNow(value) {
  const d = toDate(value)
  if (!d) {
    return ''
  }

  const diffSec = (Date.now() - d.getTime()) / 1000
  const abs = Math.abs(diffSec)
  const future = diffSec < 0

  const fmt = (num, unit) =>
    future ? `in ${num} ${unit}${num === 1 ? '' : 's'}` : `${num} ${unit}${num === 1 ? '' : 's'} ago`

  if (abs < 45) {
    return future ? 'in a few seconds' : 'a few seconds ago'
  }
  if (abs < 90) {
    return fmt(1, 'minute')
  }
  if (abs < 45 * 60) {
    return fmt(Math.round(abs / 60), 'minute')
  }
  if (abs < 90 * 60) {
    return fmt(1, 'hour')
  }
  if (abs < 24 * 3600) {
    return fmt(Math.round(abs / 3600), 'hour')
  }

  const days = Math.round(abs / 86400)
  if (days < 30) {
    return fmt(days, 'day')
  }
  if (days < 60) {
    return fmt(1, 'month')
  }
  if (days < 365) {
    return fmt(Math.round(days / 30), 'month')
  }
  return fmt(Math.round(days / 365), 'year')
}

// Whole days elapsed between now and a past date (moment diff('days')).
export function daysDiff(value) {
  const d = toDate(value)
  if (!d) {
    return Infinity
  }
  return Math.floor((Date.now() - d.getTime()) / MS_PER_DAY)
}

// Milliseconds from now until the given date (negative when past).
export function msUntil(value) {
  const d = toDate(value)
  if (!d) {
    return NaN
  }
  return d.getTime() - Date.now()
}

// Milliseconds elapsed since a past date.
export function msSince(value) {
  const d = toDate(value)
  if (!d) {
    return NaN
  }
  return Date.now() - d.getTime()
}

// Whole days elapsed since a past date, rounded up.
export function daysAgoCeil(value) {
  return Math.ceil(msSince(value) / MS_PER_DAY)
}

export function addDays(value, days) {
  const d = value == null ? new Date() : toDate(value)
  if (!d) {
    return null
  }
  const r = new Date(d)
  r.setDate(r.getDate() + days)
  return r
}

// "Jan 05"
export function formatMonthDay(value) {
  const d = toDate(value)
  if (!d) {
    return ''
  }
  return `${MONTHS[d.getMonth()]} ${pad2(d.getDate())}`
}

// "Jan 05, 2026"
export function formatMonthDayYear(value) {
  const d = toDate(value)
  if (!d) {
    return ''
  }
  return `${formatMonthDay(d)}, ${d.getFullYear()}`
}

// "January 05, 2026"
export function formatFullMonthDayYear(value) {
  const d = toDate(value)
  if (!d) {
    return ''
  }
  return `${MONTHS_FULL[d.getMonth()]} ${pad2(d.getDate())}, ${d.getFullYear()}`
}

// "Jan 05, 2026 - 3:05 PM"
export function formatMonthDayYearTime(value) {
  const d = toDate(value)
  if (!d) {
    return ''
  }
  const h24 = d.getHours()
  const ampm = h24 >= 12 ? 'PM' : 'AM'
  const h = ((h24 + 11) % 12) + 1
  return `${formatMonthDay(d)}, ${d.getFullYear()} - ${h}:${pad2(d.getMinutes())} ${ampm}`
}

// moment.calendar approximation: Today / Tomorrow / Yesterday / weekday / date.
export function calendar(value) {
  const d = value == null ? new Date() : toDate(value)
  if (!d) {
    return ''
  }

  const now = new Date()
  const startToday = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const startDay = new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
  const dayDiff = Math.round((startDay - startToday) / MS_PER_DAY)

  if (dayDiff === 0) {
    return 'Today'
  }
  if (dayDiff === 1) {
    return 'Tomorrow'
  }
  if (dayDiff === -1) {
    return 'Yesterday'
  }
  if (dayDiff > -7 && dayDiff < 7) {
    return d.toLocaleDateString('en-US', { weekday: 'long' })
  }
  return formatMonthDayYear(d)
}

export function toUnixMs(value) {
  const d = toDate(value)
  return d ? d.getTime() : NaN
}

export function dateFromNow(date) {
  const d = toDate(date)
  if (!d) {
    return ''
  }

  // Under a day: relative time. Under ~13 months: short date. Else full date.
  if (msSince(d) < MS_PER_DAY) {
    return fromNow(d)
  }
  if (msSince(d) < (1 + 30 + 365) * MS_PER_DAY) {
    return formatMonthDay(d)
  }
  return formatMonthDayYear(d)
}

export function daysFromNow(d) {
  const date = toDate(d)
  if (!date) {
    return ''
  }

  const diffDays = Math.round(msSince(date) / MS_PER_DAY)
  if (diffDays >= 20 && diffDays <= 60) {
    return `${diffDays} days ago`
  }

  return fromNow(date)
}

export function dateCalendar(date) {
  return formatFullMonthDayYear(date)
}

export function dateTime(date) {
  return formatMonthDayYearTime(date)
}

export function errorSimple(error) {
  if (!error) {
    return ''
  }

  return error.split(':')[0]
}
