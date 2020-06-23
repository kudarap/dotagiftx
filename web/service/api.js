import trimEnd from 'lodash/trimEnd'
import * as http from './http'

export const API_URL = process.env.NEXT_PUBLIC_API_URL
export const CDN_URL = `${trimEnd(process.env.NEXT_PUBLIC_CDN_URL, '/')}/`

export const fetcher = url => http.request(http.GET, url)
export const fetcherWithToken = url => http.request(http.GET, url)

// API Endpoints
const AUTH_STEAM = '/auth/steam'
const AUTH_RENEW = '/auth/renew'
const AUTH_REVOKE = '/auth/revoke'
const MY_PROFILE = '/my/profile'
const ITEMS = '/items'
const VERSION = '/'

export const authSteam = (ot, ov) =>
  http.request(http.GET, `${AUTH_STEAM}?oauth_token=${ot}&oauth_verifier=${ov}`)
export const authRenew = refreshToken =>
  http.request(http.POST, AUTH_RENEW, { refresh_token: refreshToken })
export const authRevoke = refreshToken =>
  http.request(http.POST, AUTH_REVOKE, { refresh_token: refreshToken })

export const version = () => http.request(http.GET, VERSION)

export const myProfile = {
  GET: () => http.authnRequest(http.GET, MY_PROFILE),
  PATCH: profile => http.authnRequest(http.PATCH, MY_PROFILE, profile),
}
export const item = http.baseObjectRequest(ITEMS)
export const itemSearch = http.baseSearchRequest(ITEMS)
