// Market entity constants

export const MARKET_TYPE_ASK = 10
export const MARKET_TYPE_BID = 20

export const MARKET_TYPE_MAP_TEXT = {
  [MARKET_TYPE_ASK]: 'Offer',
  [MARKET_TYPE_BID]: 'Buy Order',
}

export const MARKET_STATUS_PENDING = 100
export const MARKET_STATUS_LIVE = 200
export const MARKET_STATUS_RESERVED = 300
export const MARKET_STATUS_SOLD = 400
export const MARKET_STATUS_BID_COMPLETED = 410
export const MARKET_STATUS_REMOVED = 500
export const MARKET_STATUS_CANCELLED = 600

export const MARKET_STATUS_MAP_TEXT = {
  [MARKET_STATUS_PENDING]: 'Pending',
  [MARKET_STATUS_LIVE]: 'Listed',
  [MARKET_STATUS_RESERVED]: 'Reserved',
  [MARKET_STATUS_SOLD]: 'Delivered',
  [MARKET_STATUS_BID_COMPLETED]: 'Completed Order',
  [MARKET_STATUS_REMOVED]: 'Removed',
  [MARKET_STATUS_CANCELLED]: 'Cancelled',
}

export const MARKET_BID_STATUS_MAP_TEXT = {
  [MARKET_STATUS_LIVE]: 'Ordered',
  [MARKET_STATUS_BID_COMPLETED]: 'Order Completed',
  [MARKET_STATUS_REMOVED]: 'Order Removed',
  [MARKET_STATUS_CANCELLED]: 'Order Cancelled',
}

export const MARKET_STATUS_MAP_COLOR = {
  [MARKET_STATUS_PENDING]: 'yellow',
  [MARKET_STATUS_LIVE]: 'lightgreen',
  [MARKET_STATUS_RESERVED]: 'violet',
  [MARKET_STATUS_SOLD]: 'aqua',
  [MARKET_STATUS_BID_COMPLETED]: 'deepskyblue',
  [MARKET_STATUS_REMOVED]: 'grey',
  [MARKET_STATUS_CANCELLED]: 'orangered',
}

export const MARKET_QTY_LIMIT = 5
export const MARKET_NOTES_MAX_LEN = 200
