import MenuItem from '@mui/material/MenuItem'
import Link from '@/components/Link'
import type { Profile } from '@/lib/types'

export default function NavItems({
  profile,
  onClose,
  onLogout,
}: {
  profile: Profile
  onClose: () => void
  onLogout: () => void
}) {
  const handleClose = () => {
    onClose()
  }

  const handleLogout = () => {
    onLogout()
  }

  return [
    <MenuItem
      key="profile"
      onClick={handleClose}
      component={Link}
      nativeButton={false}
      href={`/profiles/${profile.steam_id}`}
      disableUnderline
    >
      View Profile
    </MenuItem>,
    <MenuItem
      key="listings"
      onClick={handleClose}
      component={Link}
      nativeButton={false}
      href="/my-listings"
      disableUnderline
    >
      Listings
    </MenuItem>,
    <MenuItem
      key="orders"
      onClick={handleClose}
      component={Link}
      nativeButton={false}
      href="/my-orders"
      disableUnderline
    >
      Orders
    </MenuItem>,
    <MenuItem
      key="feedback"
      onClick={handleClose}
      component={Link}
      nativeButton={false}
      href="/submit-feedback"
      disableUnderline
    >
      Feedback
    </MenuItem>,
    // <MenuItem key={key++} onClick={handleClose} component={Link} href="/updates" disableUnderline>
    //   Updates
    // </MenuItem>,
    <MenuItem key="signout" onClick={handleLogout}>
      Sign out
    </MenuItem>,
  ]
}
