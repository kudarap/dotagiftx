import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Link from '@material-ui/core/Link'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(6),
  },
}))

export default function Privacy() {
  const classes = useStyles()

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Privacy Policy</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Privacy Policy
          </Typography>
          <br />

          <Typography color="textSecondary">
            All of the data provided such as <em>steam id, profile name, and avatar image</em> comes
            from public sources such as the Steam WebAPI and Steam community profiles.
          </Typography>
          <br />

          <Typography color="textSecondary">
            We use cookies to keep your signed in session active and Google Analytics to monitor
            site traffic. For more information, see{' '}
            <Link
              href="https://www.google.com/policies/privacy/partners/"
              target="_blank"
              color="secondary"
              rel="noreferrer noopener">
              Google’s privacy terms
            </Link>
            .
          </Typography>
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
