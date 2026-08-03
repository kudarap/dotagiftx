const LATEST_TREASURE_DROP = new Date(2025, 12, 15)

const STILL_NEW_DAYS = 30

export const isTreasureNew = v => {
  const releaseDate = new Date(v)
  if (!releaseDate) {
    return false
  }

  const now = new Date()
  const diff = (now - releaseDate) / (1000 * 3600 * 24)
  return diff < STILL_NEW_DAYS
}

export const isRecentTreasureNew = () => isTreasureNew(LATEST_TREASURE_DROP)
