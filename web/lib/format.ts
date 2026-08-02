import moment from 'moment'

export function amount(n: number | string, currency = '') {
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

export function numberWithCommas(n: number | string) {
  return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

export function dateFromNow(date: string | number) {
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

export function daysFromNow(d: string | number) {
  const date = moment(d)

  const diffDays = (moment().unix() - date.unix()) / 86400
  // if (diffDays >= 20 && diffDays <= 60) {
  if (diffDays >= 20 && diffDays <= 60) {
    return `${diffDays.toFixed()} days ago`
  }

  return date.fromNow()
}

export function dateCalendar(date: string | number) {
  return moment(date).format('MMMM DD, YYYY')
}

export function dateTime(date: string | number) {
  return moment(date).format('MMM DD, YYYY - h:mm A')
}

export function errorSimple(error: string) {
  if (!error) {
    return ''
  }

  return error.split(':')[0]
}
