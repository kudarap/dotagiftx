import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import dynamic from 'next/dynamic'
import { makeStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Avatar from '@material-ui/core/Avatar'
import Toolbar from '@material-ui/core/Toolbar'
import Button from '@/components/Button'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import IconButton from '@material-ui/core/IconButton'
import MoreIcon from '@material-ui/icons/MoreVert'
import Container from '@/components/Container'
import * as Storage from '@/service/storage'
import { authRevoke, myProfile } from '@/service/api'
import { clear as destroyLoginSess } from '@/service/auth'
import Link from '@/components/Link'
import SteamIcon from '@/components/SteamIcon'
import { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import { APP_NAME } from '@/constants/strings'
import { APP_CACHE_PROFILE } from '@/constants/app'
// import SearchInputMini from '@/components/SearchInputMini'
const SearchInputMini = dynamic(() => import('@/components/SearchInputMini'))

const useStyles = makeStyles(theme => ({
  root: {},
  appBar: {
    borderTop: 'none',
    borderRight: 'none',
    borderLeft: 'none',
  },
  brand: {
    height: 30,
    marginBottom: -5,
    '-webkit-transition': 'all 1s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
    transition: 'all 1s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
    '&:hover': {
      filter: 'brightness(115%)',
    },
    // This fixes the tap highlight effect on mobile.
    '-webkit-touch-callout': 'none',
    '-webkit-user-select': 'none',
    '-khtml-user-select': 'none',
    '-moz-user-select': 'none',
    '-ms-user-select': 'none',
    'user-select': 'none',
    '-webkit-tap-highlight-color': 'transparent',
  },
  avatar: {
    width: 36,
    height: 36,
    border: `1px solid ${theme.palette.grey[700]}`,
    '&:hover': {
      borderColor: theme.palette.grey[600],
    },
    cursor: 'pointer',
  },
  avatarMenu: {
    marginTop: theme.spacing(4),
  },
  spacer: {
    width: theme.spacing(1),
  },
}))

const defaultProfile = {
  id: '',
  steam_id: '',
  name: '',
  avatar: '',
  created_at: null,
}

export default function Header({ disableSearch }) {
  const classes = useStyles()
  // NOTE! this makes the mobile version of the nav to be ignored when on homepage
  // which is the disableSearch prop uses.
  const { isMobile: isXsScreen, isLoggedIn, currentAuth } = useContext(AppContext)
  const isMobile = isXsScreen && !disableSearch

  const [profile, setProfile] = React.useState(defaultProfile)

  React.useEffect(() => {
    ;(async () => {
      let profile = Storage.get(APP_CACHE_PROFILE)
      if (profile) {
        setProfile(profile)
        return
      }

      profile = await myProfile.GET()
      Storage.save(APP_CACHE_PROFILE, profile)
      setProfile(profile)
    })()
  }, [])

  const [anchorEl, setAnchorEl] = React.useState(null)
  const handleClick = e => {
    setAnchorEl(e.currentTarget)
  }
  const handleClose = () => {
    setAnchorEl(null)
  }

  const [moreEl, setMoreEl] = React.useState(null)
  const handleMoreClick = e => {
    setMoreEl(e.currentTarget)
  }
  const handleMoreClose = () => {
    setMoreEl(null)
  }

  const handleLogout = () => {
    ;(async () => {
      try {
        await authRevoke(currentAuth.refresh_token)
      } catch (e) {
        console.warn(e.message)
      }
      destroyLoginSess()
      handleClose()
      // eslint-disable-next-line no-undef
      window.location = '/'
    })()
  }

  return (
    <AppBar position="static" variant="outlined" className={classes.appBar}>
      {/*<div style={{ textAlign: 'center', backgroundColor: 'crimson' }}>*/}
      {/*  You are viewing a development version of this site.&nbsp;*/}
      {/*  <Link href="https://dotagiftx.com">*/}
      {/*    <strong>Take me to live site</strong>*/}
      {/*  </Link>*/}
      {/*  <span style={{ float: 'right', paddingRight: 16, cursor: 'pointer' }}>close</span>*/}
      {/*</div>*/}
      <Container disableMinHeight>
        <Toolbar variant="dense" disableGutters>
          {/* Branding button */}
          {/* Desktop nav branding */}
          <Link href="/" disableUnderline>
            {!isMobile ? (
              <img
                className={classes.brand}
                src="/brand_1x.png"
                srcSet="/brand_1x.png 1x, /brand_2x.png 2x"
                alt={APP_NAME}
              />
            ) : (
              <img
                className={classes.brand}
                src="/icon_1x.png"
                srcSet="/icon_1x.png 1x, /icon_2x.png 2x"
                alt={APP_NAME}
              />
            )}
          </Link>
          <span className={classes.spacer} />
          {!disableSearch && <SearchInputMini />}

          {/* Desktop nav buttons */}
          {!isMobile && (
            <>
              <span style={{ flexGrow: 1 }} />

              {/* Post item button */}
              {/*<Button variant="outlined" component={Link} href="/buy-order" disableUnderline>*/}
              {/*  Buy Order*/}
              {/*</Button>*/}
              {/*<span className={classes.spacer} />*/}
              <Button
                variant="outlined"
                color="secondary"
                component={Link}
                href="/post-item"
                disableUnderline>
                Post Item
              </Button>
              <span className={classes.spacer} />

              {/* Avatar menu button */}
              {isLoggedIn ? (
                <>
                  <Avatar
                    aria-controls="avatar-menu"
                    aria-haspopup="true"
                    onClick={handleClick}
                    className={classes.avatar}
                    {...retinaSrcSet(profile.avatar, 36, 36)}
                  />
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
                      as={`/profiles/${profile.steam_id}`}
                      disableUnderline>
                      View Profile
                    </MenuItem>
                    <MenuItem
                      onClick={handleClose}
                      component={Link}
                      href="/my-listings"
                      disableUnderline>
                      Listings
                    </MenuItem>
                    <MenuItem
                      onClick={handleClose}
                      component={Link}
                      href="/my-reservations"
                      disableUnderline>
                      Reservations
                    </MenuItem>
                    <MenuItem
                      onClick={handleMoreClose}
                      component={Link}
                      href="/my-buyorders"
                      disableUnderline>
                      Buy Orders
                    </MenuItem>
                    <MenuItem
                      onClick={handleClose}
                      component={Link}
                      href="/my-history"
                      disableUnderline>
                      History
                    </MenuItem>
                    <MenuItem onClick={handleLogout}>Sign out</MenuItem>
                  </Menu>
                </>
              ) : (
                <Button startIcon={<SteamIcon />} component={Link} href="/login" disableUnderline>
                  Sign in
                </Button>
              )}
            </>
          )}

          {/* Mobile buttons */}
          {isMobile && (
            <>
              <span className={classes.spacer} />
              <IconButton
                aria-controls="more-menu"
                aria-haspopup="true"
                size="small"
                onClick={handleMoreClick}>
                <MoreIcon />
              </IconButton>

              <Menu
                className={classes.avatarMenu}
                id="more-menu"
                anchorEl={moreEl}
                keepMounted
                open={Boolean(moreEl)}
                onClose={handleMoreClose}>
                <MenuItem
                  onClick={handleMoreClose}
                  component={Link}
                  href="/post-item"
                  disableUnderline>
                  Post Item
                </MenuItem>

                {isLoggedIn ? (
                  [
                    <MenuItem
                      onClick={handleMoreClose}
                      component={Link}
                      href="/profiles/[id]"
                      as={`/profiles/${profile.steam_id}`}
                      disableUnderline>
                      View Profile
                    </MenuItem>,
                    <MenuItem
                      onClick={handleMoreClose}
                      component={Link}
                      href="/my-listings"
                      disableUnderline>
                      Listings
                    </MenuItem>,
                    <MenuItem
                      onClick={handleMoreClose}
                      component={Link}
                      href="/my-reservations"
                      disableUnderline>
                      Reservations
                    </MenuItem>,
                    <MenuItem
                      onClick={handleMoreClose}
                      component={Link}
                      href="/my-buyorders"
                      disableUnderline>
                      Buy Orders
                    </MenuItem>,
                    <MenuItem
                      onClick={handleMoreClose}
                      component={Link}
                      href="/my-history"
                      disableUnderline>
                      History
                    </MenuItem>,
                    <MenuItem onClick={handleLogout}>Sign out</MenuItem>,
                  ]
                ) : (
                  <MenuItem
                    onClick={handleMoreClose}
                    component={Link}
                    href="/login"
                    disableUnderline>
                    Sign in
                  </MenuItem>
                )}
              </Menu>
            </>
          )}
        </Toolbar>
      </Container>
    </AppBar>
  )
}
Header.propTypes = {
  disableSearch: PropTypes.bool,
}
Header.defaultProps = {
  disableSearch: false,
}
