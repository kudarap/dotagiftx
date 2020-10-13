import url from 'url'
import querystring from 'querystring'

export function getQuery(s) {
  return querystring.decode(url.parse(s).query)
}
