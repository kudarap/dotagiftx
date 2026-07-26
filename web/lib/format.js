import moment from 'moment'

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
  const d = moment(date)
  const dc = d.clone()
  const now = moment()

  if (now < dc.add(1, 'day')) {
    return d.fromNow()
  }
  if (now < dc.add(1, 'month')) {
    // return `${((now.unix() - d.unix()) / 86400).toFixed()} days ago`
  }
  if (now < dc.add(1, 'year')) {
    return d.format('MMM DD')
  }
  return d.format('MMM DD, YYYY')
}

export function daysFromNow(d) {
  const date = moment(d)

  const diffDays = ((moment().unix() - date.unix()) / 86400).toFixed()
  // if (diffDays >= 20 && diffDays <= 60) {
  if (diffDays >= 20 && diffDays <= 60) {
    return `${diffDays} days ago`
  }

  return date.fromNow()
}

export function dateCalendar(date) {
  return moment(date).format('MMMM DD, YYYY')
}

export function dateTime(date) {
  return moment(date).format('MMM DD, YYYY - h:mm A')
}

export function errorSimple(error) {
  if (!error) {
    return ''
  }

  return error.split(':')[0]
}
