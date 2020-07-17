import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Avatar from '@material-ui/core/Avatar'
import Toolbar from '@material-ui/core/Toolbar'
import Typography from '@material-ui/core/Typography'
import Button from '@/components/Button'
import Container from '@/components/Container'
import Link from '@/components/Link'
import SteamIcon from '@/components/SteamIcon'
import { CDN_URL, myProfile } from '@/service/api'
import { isOk as isLoggedIn } from '@/service/auth'
import * as Storage from '@/service/storage'
import SearchInputMini from '@/components/SearchInputMini'

const useStyles = makeStyles(theme => ({
  root: {},
  appBar: {
    borderTop: 'none',
    borderRight: 'none',
    borderLeft: 'none',
  },
  title: {
    [theme.breakpoints.down('sm')]: {
      fontSize: 15,
    },
    fontSize: 17,
    textShadow: '0px 0px 16px #C79123',
    // textTransform: 'uppercase',
    // fontWeight: 'bold',
    background: 'linear-gradient(#F8E8B9 10%, #fff 90%)',
    '-webkit-background-clip': 'text',
    '-webkit-text-fill-color': 'transparent',
    filter: 'drop-shadow(0px 0px 10px black)',
    letterSpacing: 2,
    cursor: 'pointer',
  },
  avatar: {
    width: theme.spacing(3),
    height: theme.spacing(3),
  },
  spacer: {
    width: theme.spacing(1),
  },
}))

const PROFILE_CACHE_KEY = 'profile'

const defaultProfile = {
  id: '',
  steam_id: '',
  name: '',
  avatar: '',
  created_at: null,
}

export default function ({ disableSearch = false }) {
  const classes = useStyles()

  const [profile, setProfile] = React.useState(defaultProfile)

  React.useEffect(() => {
    const get = async () => {
      const res = await myProfile.GET()
      setProfile(res)
      Storage.save(PROFILE_CACHE_KEY, res)
    }

    if (isLoggedIn()) {
      const hit = Storage.get(PROFILE_CACHE_KEY)
      if (hit) {
        setProfile(hit)
        return
      }
      // fetch data from api
      get()
    }
  }, [])

  return (
    <header>
      <AppBar position="static" variant="outlined" className={classes.appBar}>
        <Container disableMinHeight>
          <Toolbar variant="dense" disableGutters>
            <Link href="/" disableUnderline>
              <Typography component="h1" className={classes.title}>
                <strong>DotagiftX</strong>
              </Typography>
            </Link>
            <span className={classes.spacer} />
            {!disableSearch && <SearchInputMini />}
            <span style={{ flexGrow: 1 }} />
            <Button variant="outlined" color="secondary">
              Post Item
            </Button>
            <span className={classes.spacer} />
            {isLoggedIn() ? (
              <Button
                startIcon={
                  <Avatar className={classes.avatar} src={profile && CDN_URL + profile.avatar} />
                }>
                {profile && profile.name}
              </Button>
            ) : (
              <Button startIcon={<SteamIcon />} component={Link} href="/login">
                Sign in
              </Button>
            )}
          </Toolbar>
        </Container>
      </AppBar>
    </header>
  )
}
