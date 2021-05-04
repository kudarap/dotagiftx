import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function Blacklist() {
  const classes = useStyles()

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Guides</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Guides
          </Typography>
          <br />

          <Typography color="textSecondary">
            All of the data provided such as <em>steam id, profile name, and avatar image</em> comes
            from public sources such as the Steam WebAPI and Steam community profiles.
          </Typography>
          <br />

          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
