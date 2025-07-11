// Market entity constants

export const USER_STATUS_SUSPENDED = 300
export const USER_STATUS_BANNED = 400

export const USER_STATUS_MAP_LABEL = {
  [USER_STATUS_SUSPENDED]: 'Suspended',
  [USER_STATUS_BANNED]: 'Banned',
}

export const USER_STATUS_MAP_TEXT = {
  [USER_STATUS_SUSPENDED]: 'This account was suspended over scam report and under investigation.',
  [USER_STATUS_BANNED]: 'This account was banned over scam incident.',
}

export const USER_STATUS_MAP_COLOR = {
  [USER_STATUS_SUSPENDED]: '#aa6600',
  [USER_STATUS_BANNED]: '#a00',
}

export const USER_SUBSCRIPTION_SUPPORTER = 100
export const USER_SUBSCRIPTION_TRADER = 101
export const USER_SUBSCRIPTION_PARTNER = 109

export const USER_SUBSCRIPTION_LIST = [
  USER_SUBSCRIPTION_SUPPORTER,
  USER_SUBSCRIPTION_TRADER,
  USER_SUBSCRIPTION_PARTNER,
]

export const USER_SUBSCRIPTION_MAP_COLOR = {
  [USER_SUBSCRIPTION_SUPPORTER]: '#596b95',
  [USER_SUBSCRIPTION_TRADER]: '#629cbd',
  [USER_SUBSCRIPTION_PARTNER]: '#C79123',
}

export const USER_SUBSCRIPTION_MAP_LABEL = {
  [USER_SUBSCRIPTION_SUPPORTER]: 'Supporter',
  [USER_SUBSCRIPTION_TRADER]: 'Trader',
  [USER_SUBSCRIPTION_PARTNER]: 'Partner',
}

export const USER_AGE_CAUTION = 7

export const USER_SUBSCRIPTION_BADGE_MODE = true
