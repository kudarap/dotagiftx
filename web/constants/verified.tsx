// Verified delivery and inventory

import type { CSSProperties, ReactElement } from 'react'
// import NoHitIcon from '@mui/icons-material/RemoveCircleOutline'
// import NoHitIcon from '@mui/icons-material/HighlightOff'
import NoHitIcon from '@mui/icons-material/Block'
import CheckIcon from '@mui/icons-material/Done'
import DoubleCheckIcon from '@mui/icons-material/DoneAll'
import Private from '@mui/icons-material/VisibilityOff'
// import Private from '@mui/icons-material/Block'
import Error from '@mui/icons-material/ErrorOutlined'
import ManualCheckIcon from '@mui/icons-material/CheckCircleOutlined'
import PendingIcon from '@mui/icons-material/Pending'

const iconStyle: { style: CSSProperties } = {
  style: {
    fontSize: '1rem',
    marginLeft: 4,
    marginRight: 2,
    marginBottom: -2,
    color: 'grey',
  },
}

const rareStyle: { style: CSSProperties } = {
  style: { ...iconStyle.style, color: 'lightgreen' },
}

const resellStyle: { style: CSSProperties } = {
  style: { ...iconStyle.style, color: 'aqua' },
}

const pendingStyle: { style: CSSProperties } = {
  style: { ...iconStyle.style, color: 'grey' },
}

const ultraStyle: { style: CSSProperties } = {
  style: { ...iconStyle.style, color: 'gold' },
}

export const VERIFIED_INVENTORY_PENDING = 0
export const VERIFIED_INVENTORY_NOHIT = 100
export const VERIFIED_INVENTORY_VERIFIED = 200
export const VERIFIED_INVENTORY_VERIFIED_RESELL = 201
export const VERIFIED_INVENTORY_PRIVATE = 400
export const VERIFIED_INVENTORY_ERROR = 500

export const VERIFIED_INVENTORY_MAP_LABEL: Record<number, string> = {
  [VERIFIED_INVENTORY_NOHIT]: 'Not Found',
  [VERIFIED_INVENTORY_VERIFIED]: 'Item Verified',
  [VERIFIED_INVENTORY_PRIVATE]: 'Private Inventory',
  [VERIFIED_INVENTORY_ERROR]: 'Error',
}
export const VERIFIED_INVENTORY_MAP_TEXT: Record<number, string> = {
  [VERIFIED_INVENTORY_NOHIT]: "Item not found from seller's inventory",
  [VERIFIED_INVENTORY_VERIFIED]: "Item detected from seller's inventory",
  [VERIFIED_INVENTORY_PRIVATE]: "Seller's inventory is private",
  [VERIFIED_INVENTORY_ERROR]: 'Error processing verification',
}

export const VERIFIED_INVENTORY_MAP_ICON: Record<number, ReactElement> = {
  [VERIFIED_INVENTORY_PENDING]: <PendingIcon {...pendingStyle} />,
  [VERIFIED_INVENTORY_NOHIT]: <NoHitIcon {...iconStyle} />,
  [VERIFIED_INVENTORY_VERIFIED]: <CheckIcon {...rareStyle} />,
  [VERIFIED_INVENTORY_VERIFIED_RESELL]: <ManualCheckIcon {...resellStyle} />,
  [VERIFIED_INVENTORY_PRIVATE]: <Private {...iconStyle} />,
  [VERIFIED_INVENTORY_ERROR]: <Error {...iconStyle} />,
}

const VERIFIED_DELIVERY_NOHIT = 100
const VERIFIED_DELIVERY_NAME_VERIFIED = 200
const VERIFIED_DELIVERY_SENDER_VERIFIED = 300
const VERIFIED_DELIVERY_PRIVATE = 400
const VERIFIED_DELIVERY_ERROR = 500

export const VERIFIED_DELIVERY_MAP_LABEL: Record<number, string> = {
  [VERIFIED_DELIVERY_NOHIT]: 'Not Found',
  [VERIFIED_DELIVERY_NAME_VERIFIED]: 'Item Verified',
  [VERIFIED_DELIVERY_SENDER_VERIFIED]: 'Sender Verified',
  [VERIFIED_DELIVERY_PRIVATE]: 'Private Inventory',
  [VERIFIED_DELIVERY_ERROR]: 'Error',
}
export const VERIFIED_DELIVERY_MAP_TEXT: Record<number, string> = {
  [VERIFIED_DELIVERY_NOHIT]: "Item not found from buyer's inventory",
  [VERIFIED_DELIVERY_NAME_VERIFIED]: "Item verified from buyer's inventory",
  [VERIFIED_DELIVERY_SENDER_VERIFIED]: "Sender avatar name matched the item from buyer's inventory",
  [VERIFIED_DELIVERY_PRIVATE]: "Buyer's inventory is private",
  [VERIFIED_DELIVERY_ERROR]: 'Error processing verification',
}

export const VERIFIED_DELIVERY_MAP_ICON: Record<number, ReactElement> = {
  [VERIFIED_DELIVERY_NOHIT]: <NoHitIcon {...iconStyle} />,
  [VERIFIED_DELIVERY_NAME_VERIFIED]: <CheckIcon {...rareStyle} />,
  [VERIFIED_DELIVERY_SENDER_VERIFIED]: <DoubleCheckIcon {...ultraStyle} />,
  [VERIFIED_DELIVERY_PRIVATE]: <Private {...iconStyle} />,
  [VERIFIED_DELIVERY_ERROR]: <Error {...iconStyle} />,
}
