import React, { useContext } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import CircularProgress from '@mui/material/CircularProgress'
import Alert from '@mui/material/Alert'
import AlertTitle from '@mui/material/AlertTitle'
import { APP_NAME } from '@/constants/strings'
import { APP_CACHE_PROFILE } from '@/constants/app'
import * as Storage from '@/service/storage'
import { authSteam, getLoginURL, myProfile } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Button from '@/components/Button'
import SteamIcon from '@/components/SteamIcon'
import { set as setAuth } from '@/service/auth'
import AppContext from '@/components/AppContext'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  warningText: {
    color: theme.palette.info.main,
  },
  heading: {
    [theme.breakpoints.down('sm')]: {
      fontSize: theme.typography.h6.fontSize,
    },
  },
  list: {
    listStyle: 'none',
    '& li:before': {
      content: `'✔'`,
      marginRight: 8,
    },
    paddingLeft: theme.spacing(2),
    marginTop: 0,
  },
  banner: {
    [theme.breakpoints.down('sm')]: {
      maxWidth: 'none',
    },
    maxWidth: theme.breakpoints.values.sm,
    margin: theme.spacing(0, 0, 2, 0),
    padding: theme.spacing(1.5),
    border: '1px solid #52564e',
    background: '#2d3431',
    borderRadius: 4,
    // color: theme.palette.warning.light,
    color: '#ceb48c',
  },
}))

export default function Login() {
  const { classes } = useStyles()
  const { isLoggedIn, isMobile } = useContext(AppContext)

  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState(null)

  const router = useRouter()
  if (isLoggedIn) {
    router.push('/my-listings')
    return null
  }

  React.useEffect(() => {
    // eslint-disable-next-line no-undef
    const query = window.location.search
    const login = async () => {
      setLoading(true)
      try {
        // Store auth details.
        const auth = await authSteam(query)
        setAuth(auth)
        Storage.removeAll()

        // Store user profile.
        const profile = await myProfile.GET()
        Storage.save(APP_CACHE_PROFILE, profile)

        // eslint-disable-next-line no-undef
        window.location = '/'
      } catch (e) {
        setError(e)
        setLoading(false)
      }
    }

    if (query) {
      login()
    }
  }, [])

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Sign In</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom className={classes.heading}>
            Signing in to <strong>{APP_NAME}</strong> allows you to access additional features.
          </Typography>
          <Typography component="h2">
            <ul className={classes.list}>
              <li>Post items</li>
              <li>Track reservations</li>
              <li>Record sales history</li>
              <li>Place buy order</li>
            </ul>
          </Typography>

          <Typography className={classes.warningText}>
            This website is not affiliated with Valve Corporation or Steam.
          </Typography>
          <br />

          <Button
            fullWidth={isMobile}
            disabled={loading}
            onClick={() => setLoading(true)}
            startIcon={loading ? <CircularProgress color="secondary" size={22} /> : <SteamIcon />}
            variant="outlined"
            size="large"
            href={getLoginURL}>
            Sign in through Steam
          </Button>
          {error && <Typography color="error">{error.message}</Typography>}
          <Typography />
          <br />

          <Typography color="textSecondary" variant="body2" gutterBottom>
            By signing in, We ask for public information about your account from the{' '}
            <Link
              target="_blank"
              rel="noreferrer noopener"
              href="https://developer.valvesoftware.com/wiki/Steam_Web_API"
              color="secondary">
              Steam Web API
            </Link>{' '}
            this includes (<em>steam id, profile name, and avatar image</em>) and use cookies to
            keep your signed in session active.
          </Typography>
          <br />

          <Alert severity="warning">
            <AlertTitle>How do I know this is real?</AlertTitle>
            When you click the sign in button, you will be redirected to{' '}
            <u>https://steamcommunity.com</u> and if you are already signed into the Steam
            community, that page will allow you simply click <strong>&quot;Sign In&quot;</strong>{' '}
            without entering your password.
          </Alert>
          <Typography className={classes.banner} variant="body2" color="textSecondary" hidden>
            <strong style={{ color: 'white' }}>How do I know this is real?</strong> When you click
            the sign in button, you will be redirected to <u>https://steamcommunity.com</u> and if
            you are already signed into the Steam community, that page will allow you simply click{' '}
            <strong>&quot;Sign In&quot;</strong> without entering your password.
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
