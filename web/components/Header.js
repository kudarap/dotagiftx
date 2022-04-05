import React, { useContext, useRef, useState } from 'react'
import PropTypes from 'prop-types'
import dynamic from 'next/dynamic'
import { makeStyles } from 'tss-react/mui'
import AppBar from '@mui/material/AppBar'
import Avatar from '@/components/Avatar'
import Toolbar from '@mui/material/Toolbar'
import Button from '@mui/material/Button'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import IconButton from '@mui/material/IconButton'
import MoreIcon from '@mui/icons-material/KeyboardArrowDown'
import SearchIcon from '@mui/icons-material/Search'
import Container from '@/components/Container'
import * as Storage from '@/service/storage'
import { authRevoke, isDonationGlowExpired, myProfile } from '@/service/api'
import { clear as destroyLoginSess } from '@/service/auth'
import Link from '@/components/Link'
import SteamIcon from '@/components/SteamIcon'
import { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import { APP_NAME } from '@/constants/strings'
import { APP_CACHE_PROFILE } from '@/constants/app'
import NavItems from '@/components/NavItems'
import LatestBan from './LatestBan'
import { NoSsr, Typography } from '@mui/material'
import { usePopupState, bindHover, bindMenu } from 'material-ui-popup-state/hooks'
import HoverMenu from 'material-ui-popup-state/HoverMenu'

const SearchInputMini = dynamic(() => import('@/components/SearchInputMini'))

import brandImage from '../public/brand_2x.png'
import Image from 'next/image'

const useStyles = makeStyles()(theme => ({
  root: {},
  appBar: {
    [theme.breakpoints.down('sm')]: {
      padding: 0,
    },
    padding: theme.spacing(0, 1.5),
  },
  logo: {
    [theme.breakpoints.down('sm')]: {
      maxWidth: 30,
      overflow: 'hidden',
    },
    marginBottom: -5,
  },
  brand: {
    height: 30,
    WebkitTransition: 'all 1s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
    transition: 'all 1s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
    '&:hover': {
      filter: 'brightness(115%)',
    },
    // This fixes the tap highlight effect on mobile.
    WebkitTouchCallout: 'none',
    WebkitUserSelect: 'none',
    KhtmUserSelect: 'none',
    MozUserSelect: 'none',
    MsUserSelect: 'none',
    UserSelect: 'none',
    WebkitTapHighlightColor: 'transparent',
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
  avatarMenu: {},
  spacer: {
    width: theme.spacing(1),
  },
  nav: {
    fontWeight: theme.typography.fontWeightMedium,
    '&:hover': {
      color: '#f1e0ba',
    },
    padding: theme.spacing(0, 1.5),
    cursor: 'pointer',
  },
  navMore: {
    display: 'flex',
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
  const { classes } = useStyles()
  // NOTE! this makes the mobile version of the nav to be ignored when on homepage
  // which is the disableSearch prop uses.
  const { isTablet: isMobile, isLoggedIn, currentAuth } = useContext(AppContext)

  const [profile, setProfile] = React.useState(defaultProfile)

  // load profile data if logged in.
  React.useEffect(() => {
    ;(async () => {
      if (!isLoggedIn) {
        return
      }

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
    <AppBar position="static" variant="outlined" elevation={0} className={classes.appBar}>
      {/*<NoticeMe />*/}
      <Container disableMinHeight maxWidth="xl">
        <Toolbar variant="dense" disableGutters>
          {/* Branding button */}
          {/* Desktop nav branding */}
          <Link href="/" disableUnderline className={classes.logo}>
            <Image
              width={134}
              height={30}
              layout="fixed"
              className={classes.brand}
              src={brandImage}
              alt={APP_NAME}
            />
            {/* <img
              width={134}
              className={classes.brand}
              src="/brand_1x.png"
              srcSet="/brand_1x.png 1x, /brand_2x.png 2x"
              alt={APP_NAME}
            /> */}
          </Link>

          <span className={classes.spacer} />

          {/* <Link
            className={classes.nav}
            href="/search?q=Aghanim"
            underline="none"
            style={{
              color: '#DFE9F2',
              textShadow: '0px 0px 10px #275AF2, 2px 2px 10px #41A0F2',
            }}>
            Aghanim's 2021
          </Link> */}
          <Link className={classes.nav} href="/treasures" underline="none">
            Treasures
          </Link>
          <Link className={classes.nav} href="/plus" underline="none">
            Dotagift<span style={{ fontSize: 20 }}>+</span>
          </Link>
          <Link className={classes.nav} href="/rules" underline="none">
            Rules
          </Link>
          <Link className={classes.nav} href="/banned-users" underline="none">
            Bans
            <LatestBan />
          </Link>

          <MoreMenu />

          <NoSsr>
            <span style={{ flexGrow: 1 }} />

            <SearchInputMini style={{ width: isMobile ? '100%' : 200, marginTop: 4 }} />
            <span className={classes.spacer} />

            {/* Post item button */}
            <Button
              variant="outlined"
              color="secondary"
              component={Link}
              href="/post-item"
              disableUnderline>
              Post item
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
                  glow={isDonationGlowExpired(profile.donated_at)}
                  {...retinaSrcSet(profile.avatar, 36, 36)}
                />
                <Menu
                  className={classes.avatarMenu}
                  id="avatar-menu"
                  anchorEl={anchorEl}
                  keepMounted
                  open={Boolean(anchorEl)}
                  onClose={handleClose}>
                  <NavItems
                    profile={profile}
                    onClose={handleClose}
                    onLogout={handleLogout}
                    isMobile={isMobile}
                  />
                </Menu>
              </>
            ) : (
              <Button startIcon={<SteamIcon />} component={Link} href="/login" disableUnderline>
                Sign in
              </Button>
            )}
          </NoSsr>

          {/* Mobile nav buttons */}
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

const moreMenuLinks = [
  ['Updates', '/updates'],
  ['Guides', '/guides'],
  ['FAQs', '/faqs'],
].map(n => ({ label: n[0], path: n[1] }))

function MoreMenu() {
  const popupState = usePopupState({
    variant: 'popover',
    popupId: 'demoMenu',
  })

  const { classes } = useStyles()
  return (
    <div>
      <Typography id="more-button" className={classes.nav} {...bindHover(popupState)}>
        <div className={classes.navMore}>
          <span>More</span> <MoreIcon />
        </div>
      </Typography>
      <HoverMenu
        {...bindMenu(popupState)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
        transformOrigin={{ vertical: 'top', horizontal: 'left' }}>
        {moreMenuLinks.map(menu => (
          <MenuItem
            key={menu.path}
            onClick={popupState.close}
            component={Link}
            href={menu.path}
            disableUnderline>
            {menu.label}
          </MenuItem>
        ))}
      </HoverMenu>
    </div>
  )
}

function NoticeMe() {
  return (
    <div style={{ textAlign: 'center', backgroundColor: 'crimson' }}>
      You are viewing a development version of this site.&nbsp;
      <Link href="https://dotagiftx.com">
        <strong>Take me to live site</strong>
      </Link>
      <span style={{ float: 'right', paddingRight: 16, cursor: 'pointer' }}>close</span>
    </div>
  )
}
