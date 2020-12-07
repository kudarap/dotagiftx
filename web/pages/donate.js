import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import AwardIcon from '@material-ui/icons/Flare'
import KeyIcon from '@material-ui/icons/VpnKey'
import MoneyIcon from '@material-ui/icons/LocalAtm'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function Faq() {
  const classes = useStyles()

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Donate</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            How can I donate?
            <Typography color="textSecondary">
              Thanks for your interest supporting DotagiftX! Since its not monetizing on views and I
              do not run ads on this site(<em>fuck ads</em>).
              <br />
              Giving a like on Steam is good enough to help the site currently, but if you really
              want to donate check them out bellow.
              {/* Currently this site is running on someone else&apos;s server that is why you don&apos;t see ads */}
            </Typography>
          </Typography>
          <br />

          <div>
            <AwardIcon color="inherit" fontSize="large" />
            <Typography color="textSecondary">Give this project a Steam Award</Typography>
          </div>
          <br />

          <div>
            <KeyIcon color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Donate tradable keys and will give you a donator badge on your profile in return.
            </Typography>
          </div>
          <br />

          <div>
            <MoneyIcon color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Paypal me if you want to help future cost on server and domain.
            </Typography>
          </div>
        </Container>
      </main>

      <Footer />
    </>
  )
}
