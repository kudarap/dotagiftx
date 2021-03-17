// Verified delivery and inventory

import Tooltip from '@material-ui/core/Tooltip'
// import NoHitIcon from '@material-ui/icons/RemoveCircleOutline'
// import NoHitIcon from '@material-ui/icons/HighlightOff'
import NoHitIcon from '@material-ui/icons/Block'
import CheckIcon from '@material-ui/icons/Done'
import DoubleCheckIcon from '@material-ui/icons/DoneAll'
import Private from '@material-ui/icons/VisibilityOff'
// import Private from '@material-ui/icons/Block'
import Error from '@material-ui/icons/ErrorOutline'

const iconStyle = {
  style: {
    fontSize: '0.875rem',
    marginLeft: 4,
    marginBottom: -2,
    color: 'grey',
  },
}

const rareStyle = {
  style: { ...iconStyle.style, color: 'lightgreen' },
}

const ultraStyle = {
  style: { ...iconStyle.style, color: 'gold' },
}

function IconToolTip({ title, children }) {
  return (
    <Tooltip title={title} placement="right">
      {children}
    </Tooltip>
  )
}

const VERIFIED_INVENTORY_NOHIT = 100
const VERIFIED_INVENTORY_VERIFIED = 200
const VERIFIED_INVENTORY_PRIVATE = 400
const VERIFIED_INVENTORY_ERROR = 500

export const VERIFIED_INVENTORY_MAP_LABEL = {
  [VERIFIED_INVENTORY_NOHIT]: 'Not found',
  [VERIFIED_INVENTORY_VERIFIED]: 'Item Verified',
  [VERIFIED_INVENTORY_PRIVATE]: 'Private inventory',
  [VERIFIED_INVENTORY_ERROR]: 'Error',
}
export const VERIFIED_INVENTORY_MAP_TEXT = {
  [VERIFIED_INVENTORY_NOHIT]: "Item not found from seller's inventory",
  [VERIFIED_INVENTORY_VERIFIED]: "Item detected from seller's inventory",
  [VERIFIED_INVENTORY_PRIVATE]: "Seller's inventory is private",
  [VERIFIED_INVENTORY_ERROR]: 'Error processing verification',
}

export const VERIFIED_INVENTORY_MAP_ICON = {
  [VERIFIED_INVENTORY_NOHIT]: (
    <IconToolTip title="Item not found from seller's inventory">
      <NoHitIcon {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_INVENTORY_VERIFIED]: (
    <IconToolTip title="Verified item from seller's inventory">
      <CheckIcon {...rareStyle} />
    </IconToolTip>
  ),
  [VERIFIED_INVENTORY_PRIVATE]: (
    <IconToolTip title="Seller's inventory is private">
      <Private {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_INVENTORY_ERROR]: (
    <IconToolTip title="Error processing verification">
      <Error {...iconStyle} />
    </IconToolTip>
  ),
}

const VERIFIED_DELIVERY_NOHIT = 100
const VERIFIED_DELIVERY_NAME_VERIFIED = 200
const VERIFIED_DELIVERY_SENDER_VERIFIED = 300
const VERIFIED_DELIVERY_PRIVATE = 400
const VERIFIED_DELIVERY_ERROR = 500

export const VERIFIED_DELIVERY_MAP_TEXT = {
  [VERIFIED_DELIVERY_NOHIT]: 'Not found',
  [VERIFIED_DELIVERY_NAME_VERIFIED]: 'Item verified',
  [VERIFIED_DELIVERY_SENDER_VERIFIED]: 'Sender verified',
  [VERIFIED_DELIVERY_PRIVATE]: 'Private inventory',
  [VERIFIED_DELIVERY_ERROR]: 'Error',
}

export const VERIFIED_DELIVERY_MAP_ICON = {
  [VERIFIED_DELIVERY_NOHIT]: (
    <IconToolTip title="Item not found from buyer's inventory">
      <NoHitIcon {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_NAME_VERIFIED]: (
    <IconToolTip title="Item verified from buyer's inventory">
      <CheckIcon {...rareStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_SENDER_VERIFIED]: (
    <IconToolTip title="Sender avatar name matched the item from buyer's inventory">
      <DoubleCheckIcon {...ultraStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_PRIVATE]: (
    <IconToolTip title="Buyer's inventory is private">
      <Private {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_ERROR]: (
    <IconToolTip title="Error processing verification">
      <Error {...iconStyle} />
    </IconToolTip>
  ),
}
