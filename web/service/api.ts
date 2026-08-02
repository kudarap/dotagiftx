import querystring from 'querystring'
import trimEnd from 'lodash/trimEnd'
import * as http from './http'
import type { Auth } from './auth'

export const API_URL = process.env.NEXT_PUBLIC_API_URL as string
export const CDN_URL = trimEnd(process.env.NEXT_PUBLIC_CDN_URL as string, '/')

export const parseParams = (url: string, filter: http.QueryFilter) =>
  `${url}?${querystring.stringify(filter)}`
export const fetcher = <T>(endpoint: string, filter: http.QueryFilter = {}) =>
  http.request<T>(http.GET, parseParams(endpoint, filter))
export const fetcherBase = <T>(endpoint: string) => http.request<T>(http.GET, endpoint)
export const fetcherWithToken = <T>(endpoint: string, filter: http.QueryFilter = {}) =>
  http.authnRequest<T>(http.GET, parseParams(endpoint, filter))

// API Endpoints
const AUTH_STEAM = '/auth/steam'
const AUTH_RENEW = '/auth/renew'
const AUTH_REVOKE = '/auth/revoke'
export const MY_PROFILE = '/my/profile'
export const MY_MARKETS = '/my/markets'
export const MY_SUBSCRIPTION_PROCESS = '/my/subscription/process'
export const MY_SUBSCRIPTION_CREATE = '/my/subscription/create'
export const USERS = '/users'
export const VANITY = '/vanity'
export const ITEMS = '/items'
export const MARKETS = '/markets'
export const CATALOGS = '/catalogs'
export const CATALOGS_TREND = '/catalogs_trend'
export const STATS = '/stats'
export const STATS_TOP_ORIGINS = `${STATS}/top_origins`
export const STATS_TOP_HEROES = `${STATS}/top_heroes`
export const STATS_TOP_KEYWORDS = `${STATS}/top_keywords`
export const STATS_MARKET_SUMMARY = `${STATS}/market_summary_v2`
export const STATS_MARKET_SUMMARY_OVERALL = `${STATS}/market_summary_overall`
export const GRAPH_MARKET_SALES = `/graph/market_sales`
export const REPORTS = '/reports'
export const BLACKLIST = '/blacklists'
export const TREASURES = '/treasures'
export const HEROES = '/heroes'
const VERSION = '/'
const TRACK = '/t'

export const authSteam = (openidQuery: string) =>
  http.request<Auth>(http.GET, `${AUTH_STEAM}${openidQuery}`)
export const authRenew = (refreshToken: string) =>
  http.request<Auth>(http.POST, AUTH_RENEW, { refresh_token: refreshToken })
export const authRevoke = (refreshToken: string) =>
  http.request<Auth>(http.POST, AUTH_REVOKE, { refresh_token: refreshToken })

export const version = () => http.request(http.GET, VERSION)
export const item = (slug: string) => http.request(http.GET, `${ITEMS}/${slug}`)
export const catalog = (slug: string, marketFilter: http.QueryFilter = {}) =>
  http.request(http.GET, `${CATALOGS}/${slug}?${querystring.stringify(marketFilter)}`)
export const user = (steamID: string) => http.request(http.GET, `${USERS}/${steamID}`)
export const vanity = (vid: string) => http.request(http.GET, `${VANITY}/${vid}`)
export const statsMarketSummary = (filter: http.QueryFilter = {}) =>
  http.request(http.GET, parseParams(STATS_MARKET_SUMMARY, filter))

export const statsMarketSummaryOverall = () => http.request(http.GET, STATS_MARKET_SUMMARY_OVERALL)

export const myMarketSearch = http.baseSearchRequest(MY_MARKETS)
export const myMarket = {
  POST: (payload: unknown) => http.authnRequest(http.POST, MY_MARKETS, payload),
  PATCH: (id: string | number, payload: unknown) =>
    http.authnRequest(http.PATCH, `${MY_MARKETS}/${id}`, payload),
}
export const myProfile = {
  GET: (nocache = false) =>
    http.authnRequest(http.GET, `${MY_PROFILE}?${nocache ? 'nocache' : ''}`),
  PATCH: (profile: unknown) => http.authnRequest(http.PATCH, MY_PROFILE, profile),
}
export const createMySubscription = (planId: string) =>
  http.authnRequest(http.POST, MY_SUBSCRIPTION_CREATE, { plan_id: planId })
export const processMySubscription = (subId: string) =>
  http.authnRequest(http.POST, MY_SUBSCRIPTION_PROCESS, { subscription_id: subId })
export const reportCreate = (payload: unknown) => http.authnRequest(http.POST, REPORTS, payload)

export const itemSearch = http.baseSearchRequest(ITEMS)
export const marketSearch = http.baseSearchRequest(MARKETS)
export const catalogSearch = http.baseSearchRequest(CATALOGS)
export const catalogTrendSearch = http.baseSearchRequest(CATALOGS_TREND)
export const reportSearch = http.baseSearchRequest(REPORTS)
export const blacklistSearch = http.baseSearchRequest(BLACKLIST)

export const treasureList = () => http.request(http.GET, TREASURES)
export const getTreasure = (slug: string) => http.request(http.GET, `${TREASURES}/${slug}`)
export const heroList = () => http.request(http.GET, HEROES)
export const getHero = (id: string | number) => http.request(http.GET, `${HEROES}/${id}`)

export const trackItemViewURL = (itemID: string | number) => `${API_URL}${TRACK}?t=v&i=${itemID}`
export const trackProfileViewURL = (userID: string | number) => `${API_URL}${TRACK}?t=p&u=${userID}`
export const getLoginURL = `${API_URL}${AUTH_STEAM}`

const donationGlowExpr = 30 // days
export const isDonationGlowExpired = (donatedAt: string) => {
  if (!donatedAt) {
    return false
  }

  const d = new Date(donatedAt)
  d.setDate(d.getDate() + donationGlowExpr)
  return d > new Date()
}
