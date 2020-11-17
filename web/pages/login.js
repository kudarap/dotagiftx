import React, { useContext } from 'react'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import CircularProgress from '@material-ui/core/CircularProgress'
import { authSteam, getLoginURL } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Button from '@/components/Button'
import SteamIcon from '@/components/SteamIcon'
import { set as setAuth } from '@/service/auth'
import * as Storage from '@/service/storage'
import AppContext from '@/components/AppContext'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(6),
  },
}))

export default function Login() {
  const classes = useStyles()
  const { isLoggedIn } = useContext(AppContext)

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
        const auth = await authSteam(query)
        setAuth(auth)
        Storage.removeAll()
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
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Signing in to <strong>DotagiftX</strong> allows you to access additional features.
          </Typography>
          {/* <Typography> */}
          {/*  <ul> */}
          {/*    <li>Item Listing</li> */}
          {/*    <li>Reservation</li> */}
          {/*  </ul> */}
          {/* </Typography> */}
          <br />

          <Typography>This website is not affiliated with Valve Corporation or Steam.</Typography>
          <br />

          <Button
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
          <br />

          <Typography color="textSecondary" variant="body2">
            By signing in, We ask for public information about your account from the Steam Web API
            this includes (<em>steam id, profile name, and avatar image</em>) and use cookies to
            keep your signed in session active.
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
