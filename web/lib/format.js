import { format, formatDistanceToNow, differenceInDays, addDays, addYears } from 'date-fns'

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

  if (now < addDays(d, 1)) {
    return formatDistanceToNow(d, { addSuffix: true })
  }
  if (now < addYears(d, 1)) {
    return format(d, 'MMM dd')
  }
  return format(d, 'MMM dd, yyyy')
}

export function daysFromNow(d) {
  const date = new Date(d)

  const diffDays = ((new Date().getTime() - date.getTime()) / 86400000).toFixed()
  // if (diffDays >= 20 && diffDays <= 60) {
  if (diffDays >= 20 && diffDays <= 60) {
    return `${diffDays} days ago`
  }

  return formatDistanceToNow(date, { addSuffix: true })
}

export function fromNow(date) {
  return formatDistanceToNow(new Date(date), { addSuffix: true })
}

export function daysBetween(date) {
  return differenceInDays(new Date(), new Date(date))
}

export function dateCalendarInDays(n) {
  return format(addDays(new Date(), n), 'MM/dd/yyyy')
}

export function dateCalendar(date) {
  return format(new Date(date), 'MMMM dd, yyyy')
}

export function dateTime(date) {
  return format(new Date(date), 'MMM dd, yyyy - h:mm a')
}

export function errorSimple(error) {
  if (!error) {
    return ''
  }

  return error.split(':')[0]
}
