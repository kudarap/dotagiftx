import React from 'react'
import PropTypes from 'prop-types'
import MenuItem from '@material-ui/core/MenuItem'
import Link from '@/components/Link'

export default function MenuItems({ profile, onClose, onLogout }) {
  const handleClose = () => {
    onClose()
  }

  const handleLogout = () => {
    onLogout()
  }

  return [
    <MenuItem
      onClick={handleClose}
      component={Link}
      href="/profiles/[id]"
      as={`/profiles/${profile.steam_id}`}
      disableUnderline>
      View Profile
    </MenuItem>,
    <MenuItem onClick={handleClose} component={Link} href="/my-listings" disableUnderline>
      Listings
    </MenuItem>,
    <MenuItem onClick={handleClose} component={Link} href="/my-orders" disableUnderline>
      Orders
    </MenuItem>,
    <MenuItem onClick={handleClose} component={Link} href="/feedback" disableUnderline>
      Feedback
    </MenuItem>,
    <MenuItem onClick={handleClose} component={Link} href="/history" disableUnderline>
      Updates
    </MenuItem>,
    <MenuItem onClick={handleLogout}>Sign out</MenuItem>,
  ]
}
MenuItems.propTypes = {
  profile: PropTypes.object.isRequired,
  onClose: PropTypes.func.isRequired,
  onLogout: PropTypes.func.isRequired,
}
