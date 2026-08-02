import FormData from 'form-data'
import querystring from 'querystring'
import * as Auth from './auth'
import type { Auth as AuthData } from './auth'
import { authRenew, API_URL } from './api'

export const GET = 'GET'
export const POST = 'POST'
export const PATCH = 'PATCH'
export const DELETE = 'DELETE'

export type Method = typeof GET | typeof POST | typeof PATCH | typeof DELETE

export type QueryFilter = querystring.ParsedUrlQueryInput

interface RequestOptions {
  method?: Method
  mode: RequestMode
  headers: Record<string, string>
  body?: unknown
}

const defaultRequestOpts: RequestOptions = {
  mode: 'cors',
  headers: {
    'Content-Type': 'application/json; charSet=utf-8',
  },
}

// fetch with retry
const fetchRetry = async (url: string, options: RequestInit, n: number): Promise<Response> => {
  try {
    return await fetch(url, options)
  } catch (err) {
    if (n === 1) throw err
    return fetchRetry(url, options, n - 1)
  }
}
// default fetch retry with maximum of 3
const defaultFetchRetry = (url: string, options: RequestInit) => fetchRetry(url, options, 3)

// base http request handle json responses and internal error
const baseRequest = async <T>(method: Method, endpoint: string, body?: unknown, token?: string | null) => {
  if (!method) {
    throw Error('Request method required')
  }

  // setup request options
  const opts: RequestOptions = { ...defaultRequestOpts, method }
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

  return defaultFetchRetry(API_URL + endpoint, opts as RequestInit)
    .then(response => {
      // Catch auth error to force logout.
      if (response.status === 401) {
        Auth.clear()

        window.location.assign('/login')
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

      return json as T
    })
}

// http request that handles JSON payload.
export function request<T>(method: Method, endpoint: string, data?: unknown): Promise<T> {
  return baseRequest(method, endpoint, data)
}

// http request and handles authentication token.
export const authnRequest = async <T>(method: Method, endpoint: string, data?: unknown): Promise<T> => {
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
export function uploadFile<T>(endpoint: string, file: File): Promise<T> {
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
export function baseObjectRequest<T>(endpoint: string) {
  return {
    [GET]: (id: string | number) => authnRequest<T>(GET, `${endpoint}/${id}`),
    [POST]: (obj: unknown) => authnRequest<T>(POST, endpoint, obj),
    [PATCH]: (id: string | number, obj: unknown) => authnRequest<T>(PATCH, `${endpoint}/${id}`, obj),
    [DELETE]: (id: string | number) => authnRequest<T>(DELETE, `${endpoint}/${id}`),
  }
}

// Basic domain search request.
export function baseSearchRequest<T>(endpoint: string) {
  return (filter: QueryFilter = {}) =>
    authnRequest<T>(GET, `${endpoint}?${querystring.stringify(filter)}`)
}
