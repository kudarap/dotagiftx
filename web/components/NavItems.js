import React from 'react'
import PropTypes from 'prop-types'
import MenuItem from '@material-ui/core/MenuItem'
import Link from '@/components/Link'

let key = 100

export default function NavItems({ profile, onClose, onLogout }) {
  const handleClose = () => {
    onClose()
  }

  const handleLogout = () => {
    onLogout()
  }

  return [
    <MenuItem
      key={key++}
      onClick={handleClose}
      component={Link}
      href={`/profiles/${profile.steam_id}`}
      disableUnderline>
      View Profile
    </MenuItem>,
    <MenuItem
      key={key++}
      onClick={handleClose}
      component={Link}
      href="/my-listings"
      disableUnderline>
      Listings
    </MenuItem>,
    <MenuItem key={key++} onClick={handleClose} component={Link} href="/my-orders" disableUnderline>
      Orders
    </MenuItem>,
    <MenuItem key={key++} onClick={handleClose} component={Link} href="/feedback" disableUnderline>
      Feedback
    </MenuItem>,
    // <MenuItem key={key++} onClick={handleClose} component={Link} href="/updates" disableUnderline>
    //   Updates
    // </MenuItem>,
    <MenuItem key={key++} onClick={handleLogout}>
      Sign out
    </MenuItem>,
  ]
}
NavItems.propTypes = {
  profile: PropTypes.object.isRequired,
  onClose: PropTypes.func.isRequired,
  onLogout: PropTypes.func.isRequired,
}
