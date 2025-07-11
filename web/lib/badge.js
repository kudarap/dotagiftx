const badgePrefix = '_BADGE'

export const getUserBadgeFromBoons = (boons = []) => {
  for (const i in boons) {
    const boon = String(boons[i])
    if (boon.endsWith(badgePrefix)) {
      return boon.replace(badgePrefix, '').toLowerCase()
    }
  }
  return null
}

const tagPrefix = '_TAG'

export const getUserTagFromBoons = (boons = []) => {
  for (const i in boons) {
    const boon = String(boons[i])
    if (boon.endsWith(tagPrefix)) {
      return boon.replace(tagPrefix, '').toLowerCase()
    }
  }
  return null
}
