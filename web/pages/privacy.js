import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function Privacy() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>Privacy Policy :: {APP_NAME}</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Privacy Policy
            <Typography variant="body2" color="textSecondary">
              August 6, 2026
            </Typography>
          </Typography>
          <br />

          <Typography color="textSecondary">
            All of the data provided such as <em>steam id, profile name, and avatar image</em> comes
            from public sources such as the Steam WebAPI and Steam community profiles.
          </Typography>
          <br />

          <Typography color="textSecondary">
            We use cookies to keep your signed in session active.
          </Typography>
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
