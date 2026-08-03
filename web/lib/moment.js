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

const MONTHS_SHORT = MONTHS_FULL.map(m => m.slice(0, 3))

const DAYS_FULL = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']

const DAYS_SHORT = DAYS_FULL.map(d => d.slice(0, 3))

const THRESHOLDS = {
  s: 44,
  m: 45,
  h: 22,
  d: 26,
  w: null,
  M: 11,
}

const RELATIVE_TIME = {
  future: 'in %s',
  past: '%s ago',
  s: 'a few seconds',
  ss: '%d seconds',
  m: 'a minute',
  mm: '%d minutes',
  h: 'an hour',
  hh: '%d hours',
  d: 'a day',
  dd: '%d days',
  M: 'a month',
  MM: '%d months',
  y: 'a year',
  yy: '%d years',
}

const CALENDAR_TIME = {
  sameDay: '[Today at] LT',
  nextDay: '[Tomorrow at] LT',
  nextWeek: 'dddd [at] LT',
  lastDay: '[Yesterday at] LT',
  lastWeek: '[Last] dddd [at] LT',
  sameElse: 'L',
}

const DAYS_PER_MONTH = 365.25 / 12

function parse(value) {
  if (value instanceof Moment) {
    return new Date(value._d)
  }
  if (value === undefined || value === null) {
    return new Date()
  }
  if (value instanceof Date) {
    return new Date(value.getTime())
  }
  if (typeof value === 'number') {
    return new Date(value)
  }
  if (typeof value === 'string') {
    return new Date(value)
  }
  return new Date(NaN)
}

const pad2 = n => String(n).padStart(2, '0')

class Moment {
  constructor(value) {
    this._d = parse(value)
  }

  clone() {
    return new Moment(this._d)
  }

  unix() {
    return Math.floor(this._d.getTime() / 1000)
  }

  diff(other, units) {
    const ms = this._d.getTime() - moment(other)._d.getTime()
    if (!units) {
      return ms
    }
    switch (units) {
      case 'second':
      case 'seconds':
        return Math.trunc(ms / 1000)
      case 'minute':
      case 'minutes':
        return Math.trunc(ms / 60000)
      case 'hour':
      case 'hours':
        return Math.trunc(ms / 3600000)
      case 'day':
      case 'days':
        return Math.trunc(ms / 86400000)
      case 'week':
      case 'weeks':
        return Math.trunc(ms / 604800000)
      case 'month':
      case 'months':
        return Math.trunc(ms / (86400000 * DAYS_PER_MONTH))
      case 'year':
      case 'years':
        return Math.trunc(ms / (86400000 * 365.25))
      default:
        return ms
    }
  }

  add(n, units) {
    const d = new Date(this._d)
    switch (units) {
      case 'second':
      case 'seconds':
        d.setSeconds(d.getSeconds() + n)
        break
      case 'minute':
      case 'minutes':
        d.setMinutes(d.getMinutes() + n)
        break
      case 'hour':
      case 'hours':
        d.setHours(d.getHours() + n)
        break
      case 'day':
      case 'days':
        d.setDate(d.getDate() + n)
        break
      case 'week':
      case 'weeks':
        d.setDate(d.getDate() + n * 7)
        break
      case 'month':
      case 'months':
        d.setMonth(d.getMonth() + n)
        break
      case 'year':
      case 'years':
        d.setFullYear(d.getFullYear() + n)
        break
    }
    return new Moment(d)
  }

  from(other) {
    const ms = this._d.getTime() - moment(other)._d.getTime()
    const abs = Math.abs(ms)
    const seconds = Math.round(abs / 1000)
    const minutes = Math.round(seconds / 60)
    const hours = Math.round(minutes / 60)
    const days = Math.round(hours / 24)
    const months = Math.round(days / DAYS_PER_MONTH)
    const years = Math.round(months / 12)

    let input =
      (seconds < THRESHOLDS.s && ['s', seconds]) ||
      (minutes <= 1 && ['m']) ||
      (minutes < THRESHOLDS.m && ['mm', minutes]) ||
      (hours <= 1 && ['h']) ||
      (hours < THRESHOLDS.h && ['hh', hours]) ||
      (days <= 1 && ['d']) ||
      (days < THRESHOLDS.d && ['dd', days])

    input = input ||
      (months <= 1 && ['M']) ||
      (months < THRESHOLDS.M && ['MM', months]) ||
      (years <= 1 && ['y']) || ['yy', years]

    const key = input[0]
    const n = input[1] || 1
    const text = RELATIVE_TIME[key].replace('%d', String(n))
    return ms > 0 ? `in ${text}` : `${text} ago`
  }

  fromNow() {
    return this.from(new Moment())
  }

  calendar() {
    const today = new Moment()
    const startOfDay = new Date(today._d)
    startOfDay.setHours(0, 0, 0, 0)

    const daysAhead = (this._d.getTime() - startOfDay.getTime()) / 86400000
    let key
    if (daysAhead < -6) {
      key = 'sameElse'
    } else if (daysAhead < -1) {
      key = 'lastWeek'
    } else if (daysAhead < 0) {
      key = 'lastDay'
    } else if (daysAhead < 1) {
      key = 'sameDay'
    } else if (daysAhead < 2) {
      key = 'nextDay'
    } else if (daysAhead < 7) {
      key = 'nextWeek'
    } else {
      key = 'sameElse'
    }

    let fmt = CALENDAR_TIME[key]
      .replace(/\[([^\]]+)\]/g, '$1')
      .replace('LT', 'h:mm A')
      .replace('L', 'MM/DD/YYYY')
    return this.format(fmt)
  }

  format(fmt = '') {
    const d = this._d
    const h = d.getHours()
    const h12 = h % 12 || 12
    const tokens = {
      YYYY: String(d.getFullYear()),
      MMMM: MONTHS_FULL[d.getMonth()],
      MMM: MONTHS_SHORT[d.getMonth()],
      MM: pad2(d.getMonth() + 1),
      M: String(d.getMonth() + 1),
      DD: pad2(d.getDate()),
      D: String(d.getDate()),
      dddd: DAYS_FULL[d.getDay()],
      ddd: DAYS_SHORT[d.getDay()],
      HH: pad2(h),
      H: String(h),
      hh: pad2(h12),
      h: String(h12),
      mm: pad2(d.getMinutes()),
      A: h < 12 ? 'AM' : 'PM',
    }

    return fmt.replace(/YYYY|MMMM|MMM|MM|M|DD|D|dddd|ddd|HH|H|hh|h|mm|A/g, t => tokens[t] || t)
  }

  valueOf() {
    return this._d.getTime()
  }

  [Symbol.toPrimitive](hint) {
    if (hint === 'string') {
      return this.format('MM/DD/YYYY')
    }
    return this._d.getTime()
  }
}

function moment(value) {
  return new Moment(value)
}

export default moment
