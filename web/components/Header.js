import React, { useContext, useState } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from 'tss-react/mui'
import Image from 'next/image'
import AppBar from '@mui/material/AppBar'
import Avatar from '@/components/Avatar'
import Toolbar from '@mui/material/Toolbar'
import Button from '@mui/material/Button'
import MenuItem from '@mui/material/MenuItem'
import NoSsr from '@mui/material/NoSsr'
import Box from '@mui/system/Box'
import MoreIcon from '@mui/icons-material/KeyboardArrowDown'
import MenuIcon from '@mui/icons-material/Menu'
import HoverMenu from 'material-ui-popup-state/HoverMenu'
import { usePopupState, bindHover, bindMenu } from 'material-ui-popup-state/hooks'
import * as Storage from '@/service/storage'
import { authRevoke, isDonationGlowExpired, myProfile } from '@/service/api'
import { clear as destroyLoginSess } from '@/service/auth'
import { APP_CACHE_PROFILE } from '@/constants/app'
import Container from '@/components/Container'
import Link from '@/components/Link'
import SteamIcon from '@/components/SteamIcon'
import { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import { APP_NAME } from '@/constants/strings'
import NavItems from '@/components/NavItems'
import LatestBan from './LatestBan'
import brandImage from '../public/brand_2x.png'
import SearchDialog from './SearchDialog'
import SearchButton from './SearchButton'
import MenuDrawer from './MenuDrawer'

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
    [theme.breakpoints.down('md')]: {
      display: 'none',
    },
    fontWeight: theme.typography.fontWeightMedium,
    '&:hover': {
      color: '#f1e0ba',
    },
    padding: theme.spacing(0, 1.5),
    cursor: 'pointer',
  },
}))

const defaultProfile = {
  id: '',
  steam_id: '',
  name: '',
  avatar: '',
  created_at: null,
}

export default function Header() {
  const { isLoggedIn, currentAuth } = useContext(AppContext)

  // load profile data if logged in.
  const [profile, setProfile] = React.useState(defaultProfile)
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

  const [openDrawer, setOpenDrawer] = useState(false)
  const [openSearchDialog, setOpenSearchDialog] = useState(false)

  const handleLogout = () => {
    ;(async () => {
      try {
        await authRevoke(currentAuth.refresh_token)
      } catch (e) {
        console.warn(e.message)
      }
      destroyLoginSess()
      // eslint-disable-next-line no-undef
      window.location = '/'
    })()
  }

  const { classes } = useStyles()

  return (
    <>
      <AppBar position="static" variant="outlined" elevation={0} className={classes.appBar}>
        {/*<NoticeMe />*/}
        <Container disableMinHeight maxWidth="xl">
          <Toolbar variant="dense" disableGutters>
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
              Dotagift<span style={{ fontSize: 18, color: '#CA9039' }}>+</span>
            </Link>
            <Link className={classes.nav} href="/rules" underline="none">
              Rules
            </Link>
            <Link className={classes.nav} href="/bans" underline="none">
              Bans
              <LatestBan />
            </Link>
            <Link
              className={classes.nav}
              href="https://discord.gg/UFt9Ny42kM"
              target="_blank"
              rel="noreferrer noopener"
              underline="none">
              Discord
            </Link>
            <MoreMenu />

            <NoSsr>
              <span style={{ flexGrow: 1 }} />

              <SearchButton
                style={{ width: 180, marginTop: 4 }}
                onClick={() => setOpenSearchDialog(true)}
              />
              <span className={classes.spacer} />

              {/* Post item button */}
              <Button
                sx={{
                  display: {
                    xs: 'none',
                    md: 'inherit',
                  },
                }}
                variant="outlined"
                color="secondary"
                component={Link}
                href="/post-item"
                disableUnderline>
                Post item
              </Button>
              <Button
                onClick={() => setOpenDrawer(true)}
                variant="outlined"
                sx={{
                  display: {
                    width: 36,
                    height: 36,
                    xs: 'inherit',
                    md: 'none',
                  },
                }}>
                <MenuIcon fontSize="small" />
              </Button>
              <span className={classes.spacer} />

              {/* Avatar menu button */}
              {isLoggedIn ? (
                <AvatarMenu profile={profile} onLogout={handleLogout} />
              ) : (
                <Button
                  sx={{
                    display: {
                      xs: 'none',
                      md: 'inherit',
                    },
                  }}
                  startIcon={<SteamIcon />}
                  component={Link}
                  href="/login"
                  disableUnderline>
                  Sign in
                </Button>
              )}
            </NoSsr>

            {/* Mobile nav buttons */}
          </Toolbar>
        </Container>

        <SearchDialog open={openSearchDialog} onClose={() => setOpenSearchDialog(false)} />
        <MenuDrawer open={openDrawer} onClose={() => setOpenDrawer(false)} profile={profile} />
      </AppBar>
    </>
  )
}
Header.propTypes = {
  disableSearch: PropTypes.bool,
}
Header.defaultProps = {
  disableSearch: false,
}

function AvatarMenu({ profile, onLogout }) {
  const popupState = usePopupState({
    variant: 'popover',
    popupId: 'avatar-menu',
  })

  const { classes } = useStyles()
  return (
    <>
      <Avatar
        className={classes.avatar}
        glow={isDonationGlowExpired(profile.donated_at)}
        {...retinaSrcSet(profile.avatar, 36, 36)}
        {...bindHover(popupState)}
      />
      <HoverMenu
        className={classes.avatarMenu}
        {...bindMenu(popupState)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
        transformOrigin={{ vertical: 'top', horizontal: 'left' }}>
        <NavItems profile={profile} onClose={popupState.close} onLogout={onLogout} />
      </HoverMenu>
    </>
  )
}
AvatarMenu.propTypes = {
  profile: PropTypes.object,
  onLogout: PropTypes.func,
}

const moreMenuLinks = [
  ['Guides', '/guides'],
  ['FAQs', '/faqs'],
  ['Updates', '/updates'],
  ['Middleman', '/middlemen'],
].map(n => ({ label: n[0], path: n[1] }))

function MoreMenu() {
  const popupState = usePopupState({
    variant: 'popover',
    popupId: 'more-menu',
  })

  const { classes } = useStyles()
  return (
    <div>
      <Box
        sx={{
          display: {
            sm: 'none',
            md: 'flex',
          },
        }}
        className={classes.nav}
        {...bindHover(popupState)}>
        <span>More</span> <MoreIcon />
      </Box>
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

        {/* <MenuItem
          onClick={popupState.close}
          component={Link}
          href="https://discord.gg/UFt9Ny42kM"
          target="_blank"
          rel="noreferrer noopener"
          disableUnderline>
          Discord
        </MenuItem> */}
        {/* <MenuItem onClick={popupState.close} component={Link} href="/plus" disableUnderline>
          Dotagift<span style={{ fontSize: 20 }}>+</span>
        </MenuItem> */}
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
