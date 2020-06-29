import querystring from 'querystring'
import trimEnd from 'lodash/trimEnd'
import * as http from './http'

export const API_URL = process.env.NEXT_PUBLIC_API_URL
export const CDN_URL = `${trimEnd(process.env.NEXT_PUBLIC_CDN_URL, '/')}/`

export const fetcher2 = url => http.request(http.GET, url)

const parseParams = (url, filter) => `${url}?${querystring.stringify(filter)}`
export const fetcher = (endpoint, filter) => http.request(http.GET, parseParams(endpoint, filter))
export const fetcherWithToken = url => http.request(http.GET, url)

// API Endpoints
const AUTH_STEAM = '/auth/steam'
const AUTH_RENEW = '/auth/renew'
const AUTH_REVOKE = '/auth/revoke'
export const MY_PROFILE = '/my/profile'
export const USERS = '/users'
export const ITEMS = '/items'
export const MARKETS = '/markets'
export const MARKET_INDEX = '/market_index'
const VERSION = '/'

export const authSteam = (ot, ov) =>
  http.request(http.GET, `${AUTH_STEAM}?oauth_token=${ot}&oauth_verifier=${ov}`)
export const authRenew = refreshToken =>
  http.request(http.POST, AUTH_RENEW, { refresh_token: refreshToken })
export const authRevoke = refreshToken =>
  http.request(http.POST, AUTH_REVOKE, { refresh_token: refreshToken })

export const version = () => http.request(http.GET, VERSION)
export const item = slug => http.request(http.GET, `${ITEMS}/${slug}`)
export const user = steamID => http.request(http.GET, `${USERS}/${steamID}`)

export const myProfile = {
  GET: () => http.authnRequest(http.GET, MY_PROFILE),
  PATCH: profile => http.authnRequest(http.PATCH, MY_PROFILE, profile),
}
export const itemSearch = http.baseSearchRequest(ITEMS)