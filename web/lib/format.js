import { fromNow, format, unix } from '@/lib/date'

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
  const d = new Date(date)
  const now = new Date()

  const plusOneDay = new Date(d)
  plusOneDay.setDate(plusOneDay.getDate() + 1)
  if (now < plusOneDay) {
    return fromNow(date)
  }

  const plusOneMonth = new Date(plusOneDay)
  plusOneMonth.setMonth(plusOneMonth.getMonth() + 1)

  const plusOneYear = new Date(plusOneMonth)
  plusOneYear.setFullYear(plusOneYear.getFullYear() + 1)
  if (now < plusOneYear) {
    return format(date, 'MMM DD')
  }
  return format(date, 'MMM DD, YYYY')
}

export function daysFromNow(d) {
  const diffDays = Math.round((unix(new Date()) - unix(d)) / 86400)

  // if (diffDays >= 20 && diffDays <= 60) {
  if (diffDays >= 20 && diffDays <= 60) {
    return `${diffDays} days ago`
  }

  return fromNow(d)
}

export function dateCalendar(date) {
  return format(date, 'MMMM DD, YYYY')
}

export function dateTime(date) {
  return format(date, 'MMM DD, YYYY - h:mm A')
}

export function errorSimple(error) {
  if (!error) {
    return ''
  }

  return error.split(':')[0]
}
