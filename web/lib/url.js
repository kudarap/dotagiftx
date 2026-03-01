import url from 'url'
import querystring from 'querystring'

export function getQuery(s) {
  return querystring.decode(url.parse(s).query)
}

export function isValid(s) {
  try {
    const parsed = new URL(s)
    return parsed.protocol === 'http:' || parsed.protocol === 'https:'
  } catch (_) {
    return false
  }
}
