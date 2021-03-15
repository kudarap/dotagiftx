// Verified delivery and inventory

import Tooltip from '@material-ui/core/Tooltip'
import NoHitIcon from '@material-ui/icons/Block'
import CheckIcon from '@material-ui/icons/Done'
import DoubleCheckIcon from '@material-ui/icons/DoneAll'
import Private from '@material-ui/icons/VisibilityOff'
import Error from '@material-ui/icons/Error'

const iconStyle = {
  style: {
    fontSize: '0.875rem',
    marginLeft: 4,
    marginBottom: -2,
  },
}

function IconToolTip({ title, children }) {
  return (
    <Tooltip title={title} placement="right">
      {children}
    </Tooltip>
  )
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
    <IconToolTip title="Not found">
      <NoHitIcon {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_NAME_VERIFIED]: (
    <IconToolTip title="Item verified">
      <CheckIcon {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_SENDER_VERIFIED]: (
    <IconToolTip title="Sender verified">
      <CheckIcon {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_PRIVATE]: (
    <IconToolTip title="Private inventory">
      <Private {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_DELIVERY_ERROR]: (
    <IconToolTip title="Error">
      <Error {...iconStyle} />
    </IconToolTip>
  ),
}

const VERIFIED_INVENTORY_NOHIT = 100
const VERIFIED_INVENTORY_VERIFIED = 200
const VERIFIED_INVENTORY_PRIVATE = 400
const VERIFIED_INVENTORY_ERROR = 500

export const VERIFIED_INVENTORY_MAP_TEXT = {
  [VERIFIED_INVENTORY_NOHIT]: 'Not found',
  [VERIFIED_INVENTORY_VERIFIED]: 'Verified',
  [VERIFIED_INVENTORY_PRIVATE]: 'Private inventory',
  [VERIFIED_INVENTORY_ERROR]: 'Error',
}

export const VERIFIED_INVENTORY_MAP_ICON = {
  [VERIFIED_INVENTORY_NOHIT]: (
    <IconToolTip title="Not found">
      <NoHitIcon {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_INVENTORY_VERIFIED]: (
    <IconToolTip title="Verified">
      <CheckIcon {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_INVENTORY_PRIVATE]: (
    <IconToolTip title="Private inventory">
      <Private {...iconStyle} />
    </IconToolTip>
  ),
  [VERIFIED_INVENTORY_ERROR]: (
    <IconToolTip title="Error">
      <Error {...iconStyle} />
    </IconToolTip>
  ),
}
