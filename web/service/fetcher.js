import querystring from 'querystring'
import * as http from './http'

function parseParams(url, filter) {
  return `/${url}?${querystring.stringify(filter)}`
}

export const fetcher = url => http.request(http.GET, url)
// export const fetcher = (url, filter) => http.request(http.GET, parseParams(url, filter))

export const fetcherWithToken = (url, filter) => http.request(http.GET, parseParams(url, filter))
