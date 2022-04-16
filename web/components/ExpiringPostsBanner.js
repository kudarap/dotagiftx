import React, { useState } from 'react'
import PropTypes from 'prop-types'
import Alert from '@mui/material/Alert'
import IconButton from '@mui/material/IconButton'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import MoreVertIcon from '@mui/icons-material/MoreVert'
import Link from './Link'
import useLocalStorage from './useLocalStorage'

const targetUpdateID = 20220415

function ExpiringPostsBanner({ userID }) {
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
      Major update: We will role out Expiring items on May 1, 2022 —{' '}
      <Link href="/expiring-posts">Read more</Link>
    </Alert>
  )
}
ExpiringPostsBanner.propTypes = {
  userID: PropTypes.string,
}
ExpiringPostsBanner.defaultProps = {
  userID: '',
}

function BasicMenu({ onClose }) {
  const [anchorEl, setAnchorEl] = React.useState(null)
  const open = Boolean(anchorEl)
  const handleClick = event => {
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
