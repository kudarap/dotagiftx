// LocalStorage with cache mechanism.

const CACHE_KEY = 'cache'

const hash = (value: unknown) => {
  const str = JSON.stringify(value) ?? ''
  let hash = 0
  for (let i = 0, len = str.length; i < len; i++) {
    const chr = str.charCodeAt(i)
    hash = (hash << 5) - hash + chr
    hash |= 0 // Convert to 32bit integer
  }
  return hash
}

const keyPrefix = (key: string) => `${CACHE_KEY}:${String(key).split('/').shift()}`

const cKey = (key: string) => `${keyPrefix(key)}:${hash(key)}`

const now = () => new Date().getTime()

const isExpired = (ttl: number | null) => {
  // Immortal entry do not delete.
  if (ttl === null) {
    return false
  }

  return ttl < now()
}

const matchKeys = (prefix: string) => {
  const keys = []

  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i)
    if (key === null || !key.startsWith(prefix)) {
      continue
    }

    keys.push(key)
  }

  return keys
}

// Checks for expired items and remove them.
const sweep = () => {
  matchKeys(keyPrefix(CACHE_KEY)).forEach(key => {
    const { ttl } = JSON.parse(localStorage.getItem(key) ?? 'null')
    if (!isExpired(ttl)) {
      return
    }

    localStorage.removeItem(key)
  })
}

// remove by exact key.
export const remove = (key: string) => {
  localStorage.removeItem(cKey(key))
}

// remove entries with matched prefix key.
export const removeAll = (key: string = '') => {
  matchKeys(keyPrefix(key || '')).forEach(k => localStorage.removeItem(k))
}

export const get = (key: string) => {
  const raw = localStorage.getItem(cKey(key))
  if (raw === null) {
    return null
  }

  const { data, ttl } = JSON.parse(raw)
  // Return non expiry item.
  if (ttl === null) {
    return data
  }

  // Remove expired item.
  if (isExpired(ttl)) {
    remove(key)
    return null
  }

  return data
}

export const save = (key: string, data: unknown, sec: number | null = null) => {
  // Free up expired items.
  sweep()

  // Skip saving null data
  if (data === null) {
    return
  }

  let ttl = sec
  if (sec !== null) {
    // Converts TTL seconds to milli sec.
    ttl = Number(sec) * 1000
    // and adds now milli sec.
    ttl += now()
  }

  const item = { data, ttl }
  localStorage.setItem(cKey(key), JSON.stringify(item))
}
