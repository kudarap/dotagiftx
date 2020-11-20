import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Button from '@/components/Button'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import Link from '@/components/Link'
import { retinaSrcSet } from '@/components/ItemImage'

const useStyles = makeStyles(theme => ({
  avatar: {
    width: theme.spacing(3),
    height: theme.spacing(3),
  },
  avatarMenu: {
    marginTop: theme.spacing(4),
  },
}))

export default function AvatarMenu({ profile }) {
  const classes = useStyles()

  const [anchorEl, setAnchorEl] = React.useState(null)

  const handleClick = event => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  return (
    <>
      <Button
        aria-controls="avatar-menu"
        aria-haspopup="true"
        onClick={handleClick}
        startIcon={<Avatar className={classes.avatar} {...retinaSrcSet(profile.avatar, 24, 24)} />}>
        {profile && profile.name}
      </Button>
      <Menu
        className={classes.avatarMenu}
        id="avatar-menu"
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}>
        <MenuItem
          onClick={handleClose}
          component={Link}
          href="/profiles/[id]"
          as={`/profiles/${profile.steam_id}`}>
          Profile
        </MenuItem>
        <MenuItem onClick={handleClose}>Listings</MenuItem>
        {/* <MenuItem onClick={handleClose}>Buy Orders</MenuItem> */}
        <MenuItem onClick={handleClose}>Sign out</MenuItem>
      </Menu>
    </>
  )
}
AvatarMenu.propTypes = {
  profile: PropTypes.object.isRequired,
}
