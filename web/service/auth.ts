import Cookies from 'js-cookie'
import moment from 'moment'

const AUTH_KEY = 'dgAu'
const userIDKey = 'user_id'

export interface Auth {
  user_id?: number
  token?: string | null
  refresh_token?: string
  expires_at?: string
}

export const get = (): Auth => {
  const raw = Cookies.get(AUTH_KEY)
  if (!raw) {
    return {}
  }
  return JSON.parse(raw) as Auth
}

export const isOk = () => {
  return get().hasOwnProperty(userIDKey)
}

export const set = (data: Auth) => {
  Cookies.set(AUTH_KEY, JSON.stringify(data), { expires: 365, secure: true, sameSite: 'strict' })
}

export const clear = () => {
  Cookies.remove(AUTH_KEY)
}

export function getAccessToken(): string | null {
  if (!isOk()) {
    return null
  }

  const auth = get()
  return auth.token || null
}

const renewLeeway = 60 // seconds before expiration

export function isAccessTokenExpired(): boolean {
  const auth = get()
  return moment(auth.expires_at).diff(moment()) <= renewLeeway
}

const authService = {
  isOk,
  set,
  get,
  clear,
}

export default authService
