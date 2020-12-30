import Cookies from 'js-cookie'
import moment from 'moment'
// import { authRenew } from './api'

const AUTH_KEY = 'dgxAu'
const userIDKey = 'user_id'

export const get = () => {
  return Cookies.getJSON(AUTH_KEY) || {}
}

export const isOk = () => {
  // eslint-disable-next-line no-prototype-builtins
  return get().hasOwnProperty(userIDKey)
}

export const set = data => {
  let opts = null
  if (navigator.userAgent.indexOf('Safari') === -1) {
    opts = { expires: 30, secure: true, sameSite: 'strict' }
  }

  Cookies.set(AUTH_KEY, data, opts)
}

export const clear = () => {
  Cookies.remove(AUTH_KEY)
}

export function getAccessToken() {
  if (!isOk()) {
    return null
  }

  const auth = get()
  return auth.token || null
}

const renewLeeway = 60 // seconds before expiration

export function isAccessTokenExpired() {
  const auth = get()
  return moment(auth.expires_at).diff(moment()) <= renewLeeway
}

// export function renewAccessToken(onSuccess = () => {}, onError = () => {}) {
//   const auth = get()
//   // check expired access token
//   if (moment(auth.expires_at).diff(moment()) >= renewLeeway) {
//     return
//   }
//
//   // renew access token using refresh token and save
//   authRenew(auth.refresh_token)
//     .then(res => {
//       auth.token = res.token
//       auth.expires_at = res.expires_at
//       set(auth)
//       onSuccess(auth)
//     })
//     .catch(e => {
//       onError(e)
//     })
// }

export default {
  isOk,
  set,
  get,
  clear,
}
