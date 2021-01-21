import moment from 'moment'

export function amount(n, currency = '') {
  let sign = ''
  if (currency) {
    // eslint-disable-next-line default-case
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

  console.log('isDateWithin20to70days?', isDateWithin20to70days(d))

  // formats 30-60 days as days ago.
  // if (now.clone().add(-40, 'day') >= date && now.clone().add(1, 'month') >= date) {
  //   const days = ((now.unix() - date.unix()) / 86400).toFixed()
  //   return `${days} days ago`
  // }
  return date.fromNow()
}

export function dateCalendar(date) {
  return moment(date).format('MMMM DD, YYYY')
}

export function errorSimple(error) {
  if (!error) {
    return ''
  }

  return error.split(':')[0]
}

function isDateWithin20to70days(d) {
  const date = moment(d)
  const now = moment()

  console.log('-------------------------------------------')

  console.log('date', date.calendar())
  console.log('now from', now.clone().add(-40, 'days').calendar())
  console.log('now to', now.clone().add(1, 'month').calendar())
  // console.log(now.calendar())

  return false
}
