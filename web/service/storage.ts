// LocalStorage with cache mechanism.

const CACHE_KEY = 'cache'

const hash = (str: unknown): number => {
  const s = JSON.stringify(str) ?? ''
  let h = 0
  for (let i = 0, len = s.length; i < len; i++) {
    const chr = s.charCodeAt(i)
    h = (h << 5) - h + chr
    h |= 0 // Convert to 32bit integer
  }
  return h
}

const keyPrefix = (key: string) => `${CACHE_KEY}:${String(key).split('/').shift()}`

const cKey = (key: string) => `${keyPrefix(key)}:${hash(key)}`

const now = () => new Date().getTime()

const isExpired = (ttl: number | null): boolean => {
  // Immortal entry do not delete.
  if (ttl === null) {
    return false
  }

  return ttl < now()
}

const matchKeys = (prefix: string): string[] => {
  const keys: string[] = []

  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i) ?? ''
    if (!key.startsWith(prefix)) {
      continue
    }

    keys.push(key)
  }

  return keys
}

interface CacheItem<T> {
  data: T
  ttl: number | null
}

// Checks for expired items and remove them.
const sweep = () => {
  matchKeys(keyPrefix(CACHE_KEY)).forEach(key => {
    const { ttl } = JSON.parse(localStorage.getItem(key) ?? 'null') as CacheItem<unknown>
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
export const removeAll = (key = '') => {
  matchKeys(keyPrefix(key)).forEach(k => localStorage.removeItem(k))
}

export const get = <T>(key: string): T | null => {
  const item = JSON.parse(localStorage.getItem(cKey(key)) ?? 'null') as CacheItem<T> | null
  if (item === null) {
    return null
  }

  const { data, ttl } = item
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

export const save = <T>(key: string, data: T, sec: number | null = null) => {
  // Free up expired items.
  sweep()

  // Skip saving null data
  if (data === null) {
    return
  }

  let ttl: number | null = sec
  if (sec !== null) {
    // Converts TTL seconds to milli sec.
    ttl = Number(sec) * 1000
    // and adds now milli sec.
    ttl += now()
  }

  const item: CacheItem<T> = { data, ttl }
  localStorage.setItem(cKey(key), JSON.stringify(item))
}
