import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Button from '@/components/Button'
import SteamIcon from '@/components/SteamIcon'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(10),
  },
}))

export default function Login() {
  const classes = useStyles()

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
          <Button startIcon={<SteamIcon />} variant="outlined" size="large">
            Sign in through Steam
          </Button>
          <br />
          <br />
          <Typography>We use cookies to keep your signed in session active.</Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
