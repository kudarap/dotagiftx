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
