// LocalStorage with cache mechanism.

import crypto from 'crypto'

const CACHE_KEY = 'cache'

const hash = data => crypto.createHash('md5').update(data).digest('hex')

const keyPrefix = key => `${CACHE_KEY}:${String(key).split('/').shift()}`

const cKey = key => `${keyPrefix(key)}:${hash(key)}`

const now = () => new Date().getTime()

const isExpired = ttl => {
  // Immortal entry do not delete.
  if (ttl === null) {
    return false
  }

  return ttl < now()
}

const matchKeys = prefix => {
  const keys = []

  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i)
    if (!key.startsWith(prefix)) {
      continue
    }

    keys.push(key)
  }

  return keys
}

// Checks for expired items and remove them.
const sweep = () => {
  matchKeys(keyPrefix(CACHE_KEY)).forEach(key => {
    const { ttl } = JSON.parse(localStorage.getItem(key))
    if (!isExpired(ttl)) {
      return
    }

    localStorage.removeItem(key)
  })
}

// remove by exact key.
export const remove = key => {
  localStorage.removeItem(cKey(key))
}

// remove entries with matched prefix key.
export const removeAll = key => {
  matchKeys(keyPrefix(key || '')).forEach(k => localStorage.removeItem(k))
}

export const get = key => {
  const item = JSON.parse(localStorage.getItem(cKey(key)))
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

export const save = (key, data, sec = null) => {
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
