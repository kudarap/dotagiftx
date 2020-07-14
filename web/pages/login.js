import React from 'react'
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

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(10),
  },
}))

export default function Login() {
  const classes = useStyles()

  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState(null)

  React.useEffect(() => {
    // eslint-disable-next-line no-undef
    const query = window.location.search
    const login = async () => {
      setLoading(true)
      try {
        const auth = await authSteam(query)
        setAuth(auth)
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
          <Typography variant="h6" component="h1">
            Logging into DotagiftX allows you to access additional features.
          </Typography>
          <Typography>
            To provide a better service we will fetch public information about your account from the
            Steam Web API <em>(this includes steamid, profile name, and avatar)</em>.
          </Typography>
          <br />
          <Typography>This website is not affiliated with Valve Corporation or Steam.</Typography>
          <br />
          <Button
            disabled={loading}
            startIcon={loading ? <CircularProgress color="secondary" size={22} /> : <SteamIcon />}
            variant="outlined"
            size="large"
            href={getLoginURL}>
            Sign in through Steam
          </Button>
          {error && <Typography color="error">{error.message}</Typography>}
          <Typography />
          <br />
          <Typography>We use cookies to keep your signed in session active.</Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
