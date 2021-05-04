import querystring from 'querystring'
import trimEnd from 'lodash/trimEnd'
import * as http from './http'

export const API_URL = process.env.NEXT_PUBLIC_API_URL
export const CDN_URL = trimEnd(process.env.NEXT_PUBLIC_CDN_URL, '/')

export const parseParams = (url, filter) => `${url}?${querystring.stringify(filter)}`
export const fetcher = (endpoint, filter) => http.request(http.GET, parseParams(endpoint, filter))
export const fetcherBase = endpoint => http.request(http.GET, endpoint)
export const fetcherWithToken = (endpoint, filter) =>
  http.authnRequest(http.GET, parseParams(endpoint, filter))

// API Endpoints
const AUTH_STEAM = '/auth/steam'
const AUTH_RENEW = '/auth/renew'
const AUTH_REVOKE = '/auth/revoke'
export const MY_PROFILE = '/my/profile'
export const MY_MARKETS = '/my/markets'
export const USERS = '/users'
export const VANITY = '/vanity'
export const ITEMS = '/items'
export const MARKETS = '/markets'
export const CATALOGS = '/catalogs'
export const CATALOGS_TREND = '/catalogs_trend'
export const STATS = '/stats'
export const STATS_TOP_ORIGINS = `${STATS}/top_origins`
export const STATS_TOP_HEROES = `${STATS}/top_heroes`
export const STATS_MARKET_SUMMARY = `${STATS}/market_summary`
export const GRAPH_MARKET_SALES = `/graph/market_sales`
export const REPORTS = '/reports'
export const BLACKLIST = '/blacklists'
const VERSION = '/'
const TRACK = '/t'

export const authSteam = openidQuery => http.request(http.GET, `${AUTH_STEAM}${openidQuery}`)
export const authRenew = refreshToken =>
  http.request(http.POST, AUTH_RENEW, { refresh_token: refreshToken })
export const authRevoke = refreshToken =>
  http.request(http.POST, AUTH_REVOKE, { refresh_token: refreshToken })

export const version = () => http.request(http.GET, VERSION)
export const item = slug => http.request(http.GET, `${ITEMS}/${slug}`)
export const catalog = slug => http.request(http.GET, `${CATALOGS}/${slug}`)
export const user = steamID => http.request(http.GET, `${USERS}/${steamID}`)
export const vanity = vid => http.request(http.GET, `${VANITY}/${vid}`)
export const statsMarketSummary = (filter = {}) =>
  http.request(http.GET, parseParams(STATS_MARKET_SUMMARY, filter))

export const myMarketSearch = http.baseSearchRequest(MY_MARKETS)
export const myMarket = {
  POST: payload => http.authnRequest(http.POST, MY_MARKETS, payload),
  PATCH: (id, payload) => http.authnRequest(http.PATCH, `${MY_MARKETS}/${id}`, payload),
}
export const myProfile = {
  GET: () => http.authnRequest(http.GET, MY_PROFILE),
  PATCH: profile => http.authnRequest(http.PATCH, MY_PROFILE, profile),
}
export const reportCreate = payload => http.authnRequest(http.POST, REPORTS, payload)

export const itemSearch = http.baseSearchRequest(ITEMS)
export const marketSearch = http.baseSearchRequest(MARKETS)
export const catalogSearch = http.baseSearchRequest(CATALOGS)
export const catalogTrendSearch = http.baseSearchRequest(CATALOGS_TREND)
export const reportSearch = http.baseSearchRequest(REPORTS)

export const trackItemViewURL = itemID => `${API_URL}${TRACK}?t=v&i=${itemID}`
export const trackProfileViewURL = userID => `${API_URL}${TRACK}?t=p&u=${userID}`
export const getLoginURL = `${API_URL}${AUTH_STEAM}`
