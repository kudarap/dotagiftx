// import { fetch } from 'whatwg-fetch'
// import fetch from 'unfetch'
import FormData from 'form-data'
import querystring from 'querystring'
import * as Auth from './auth'
import { authRenew, API_URL } from './api'

export const GET = 'GET'
export const POST = 'POST'
export const PATCH = 'PATCH'
export const DELETE = 'DELETE'

const defaultRequestOpts = {
  mode: 'cors',
  // signal: controller.signal,
  headers: {
    'Content-Type': 'application/json; charSet=utf-8',
  },
}

// fetch with retry
const fetchRetry = async (url, options, n) => {
  try {
    return await fetch(url, options)
  } catch (err) {
    if (n === 1) throw err
    return fetchRetry(url, options, n - 1)
  }
}
// default fetch retry with maximum of 3
const defaultFetchRetry = (url, options) => fetchRetry(url, options, 3)

// base http request handle json responses and internal error
const baseRequest = (method, endpoint, body, token = null) => {
  if (method === '') {
    throw Error('Request method required')
  }

  // setup request options
  const opts = { ...defaultRequestOpts, method }
  // set access token when available
  if (token) {
    opts.headers.Authorization = `Bearer ${token}`
  }
  // GET request cant have body.
  if (method !== GET) {
    if (body instanceof FormData) {
      delete opts.headers['Content-Type']
      opts.body = body
    } else {
      opts.body = JSON.stringify(body)
    }
  }

  return defaultFetchRetry(API_URL + endpoint, opts)
    .then(response => {
      // Catch auth error to force logout.
      if (response.status === 401) {
        Auth.clear()

        window.location = '/login'
        throw Error('Authentication error')
      }
      // Catch internal error.
      if (response.status === 500) {
        throw Error('Something went wrong')
      }

      // Good response data.
      return response.json()
    })
    .then(json => {
      // Handle user error.
      if (json.error) {
        throw Error(json.msg || json.type || 'server error')
      }

      return json
    })
}

// http request that handles JSON payload.
export function request(method, endpoint, data) {
  return baseRequest(method, endpoint, data)
}

// http request and handles authentication token.
export const authnRequest = async (method, endpoint, data) => {
  // check and set access token.
  let auth = Auth.get()
  if (auth.refresh_token && (Auth.isAccessTokenExpired() || auth.token === null)) {
    const newAuth = await authRenew(auth.refresh_token)
    auth = { ...auth, ...newAuth }
    Auth.set(auth)
  }

  return baseRequest(method, endpoint, data, auth.token)
}

// Upload form file with authorization.
export function uploadFile(endpoint, file) {
  // Blob file handling and form data composition.
  const data = new FormData()
  if (file.constructor === Blob) {
    data.append('file', file, file.name)
  } else {
    data.append('file', file)
  }

  return authnRequest(POST, endpoint, data)
}

// Basic domain object request that supports all request method.
export function baseObjectRequest(endpoint) {
  return {
    [GET]: id => authnRequest(GET, `${endpoint}/${id}`),
    [POST]: obj => authnRequest(POST, endpoint, obj),
    [PATCH]: (id, obj) => authnRequest(PATCH, `${endpoint}/${id}`, obj),
    [DELETE]: id => authnRequest(DELETE, `${endpoint}/${id}`),
  }
}

// Basic domain search request.
export function baseSearchRequest(endpoint) {
  return (filter = {}) => authnRequest(GET, `${endpoint}?${querystring.stringify(filter)}`)
}

// HTTP Interceptor
// fetch = (originalFetch => {
//   return (...args) => {
//     console.log('before send')
//
//     const result = originalFetch(...args)
//     return result.then(resp => {
//       console.log('Request was sent', resp.status)
//       return resp
//     })
//   }
// })(fetch)

// Response timeout
// const timeout = 4000
// const controller = new AbortController()
// // eslint-disable-next-line no-unused-vars
// const timeoutId = setTimeout(() => controller.abort(), timeout)
