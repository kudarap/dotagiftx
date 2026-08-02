import React, { useState } from 'react'
import Alert from '@mui/material/Alert'
import IconButton from '@mui/material/IconButton'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import MoreVertIcon from '@mui/icons-material/MoreVert'
import Link from './Link'
import useLocalStorage from './useLocalStorage'

const targetUpdateID = 20220415

function ExpiringPostsBanner({ userID }: { userID?: string }) {
  const wuid = `whatsnew_id_${userID}`
  const [clientUpdateID, setClientUpdateID] = useLocalStorage(wuid, 0)

  const [open, setOpen] = useState(targetUpdateID > clientUpdateID)
  const handleClose = () => {
    setClientUpdateID(targetUpdateID)
    setOpen(false)
  }

  const handleSubmit = () => {
    handleClose()
  }

  if (!open) {
    return null
  }

  return (
    <Alert
      severity="warning"
      action={<BasicMenu color="inherit" size="small" onClose={handleSubmit} />}>
      <strong>Major update</strong>: We will roll out Expiring items on May 1, 2022 —{' '}
      <Link href="/expiring-posts">Read more</Link>
    </Alert>
  )
}

function BasicMenu({
  onClose,
  color,
  size,
}: {
  onClose: () => void
  color?: string
  size?: string
}) {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)
  const open = Boolean(anchorEl)
  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }
  const handleClose = () => {
    setAnchorEl(null)
  }

  return (
    <div>
      <IconButton
        id="expiring-banner-more-menu"
        size="small"
        color={color as 'inherit'}
        aria-controls={open ? 'basic-menu' : undefined}
        aria-haspopup="true"
        aria-expanded={open ? 'true' : undefined}
        onClick={handleClick}>
        <MoreVertIcon fontSize="inherit" />
      </IconButton>
      <Menu
        id="basic-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        MenuListProps={{
          'aria-labelledby': 'basic-button',
        }}>
        <MenuItem onClick={onClose}>Close</MenuItem>
      </Menu>
    </div>
  )
}

export default ExpiringPostsBanner
