const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
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
const DAYS = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
const MS_PER_DAY = 86400000

const pad = n => String(n).padStart(2, '0')

function toDate(date) {
  if (date instanceof Date) {
    return new Date(date.getTime())
  }
  return new Date(date)
}

export function addDays(date, days) {
  const d = toDate(date)
  d.setDate(d.getDate() + days)
  return d
}

export function unix(date) {
  return Math.floor(toDate(date).getTime() / 1000)
}

export function diffMs(a, b) {
  return toDate(a).getTime() - toDate(b).getTime()
}

export function diffDays(a, b) {
  return Math.trunc(diffMs(a, b) / MS_PER_DAY)
}

export function format(date, pattern) {
  const d = toDate(date)
  return pattern.replace(/MMMM|MMM|MM|M|DD|D|YYYY|hh|h|mm|A|a/g, token => {
    switch (token) {
      case 'MMMM':
        return MONTHS_FULL[d.getMonth()]
      case 'MMM':
        return MONTHS[d.getMonth()]
      case 'MM':
        return pad(d.getMonth() + 1)
      case 'M':
        return String(d.getMonth() + 1)
      case 'DD':
        return pad(d.getDate())
      case 'D':
        return String(d.getDate())
      case 'YYYY':
        return String(d.getFullYear())
      case 'hh':
        return pad(d.getHours() % 12 || 12)
      case 'h':
        return String(d.getHours() % 12 || 12)
      case 'mm':
        return pad(d.getMinutes())
      case 'A':
        return d.getHours() < 12 ? 'AM' : 'PM'
      case 'a':
        return d.getHours() < 12 ? 'am' : 'pm'
      default:
        return token
    }
  })
}

export function calendar(date) {
  const d = toDate(date)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const target = new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
  const dayDiff = Math.round((target - today) / MS_PER_DAY)
  const time = format(d, 'h:mm A')

  if (dayDiff === 0) {
    return `Today at ${time}`
  }
  if (dayDiff === 1) {
    return `Tomorrow at ${time}`
  }
  if (dayDiff === -1) {
    return `Yesterday at ${time}`
  }
  if (dayDiff > 1 && dayDiff < 7) {
    return `${DAYS[d.getDay()]} at ${time}`
  }
  if (dayDiff < -1 && dayDiff > -7) {
    return `Last ${DAYS[d.getDay()]} at ${time}`
  }
  return format(d, 'MM/DD/YYYY')
}

export function fromNow(date) {
  const seconds = (Date.now() - toDate(date).getTime()) / 1000

  if (Math.abs(seconds) < 45) {
    return 'a few seconds ago'
  }

  const minutes = seconds / 60
  if (Math.abs(minutes) < 45) {
    const m = Math.round(Math.abs(minutes))
    return m === 1 ? 'a minute ago' : `${m} minutes ago`
  }

  const hours = minutes / 60
  if (Math.abs(hours) < 22) {
    const h = Math.round(Math.abs(hours))
    return h === 1 ? 'an hour ago' : `${h} hours ago`
  }

  const days = hours / 24
  if (Math.abs(days) < 26) {
    const dd = Math.round(Math.abs(days))
    return dd === 1 ? 'a day ago' : `${dd} days ago`
  }

  const months = days / 30
  if (Math.abs(months) < 11) {
    const m = Math.round(Math.abs(months))
    return m === 1 ? 'a month ago' : `${m} months ago`
  }

  const years = days / 365
  const y = Math.round(Math.abs(years))
  return y === 1 ? 'a year ago' : `${y} years ago`
}
