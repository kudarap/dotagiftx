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
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  item: {
    display: 'flex',
    alignItems: 'center',
  },
  icon: {
    marginRight: theme.spacing(2),
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
              do not run ads on this site. Giving a{' '}
              <Link href="https://steamcommunity.com/sharedfiles/filedetails/?id=2313234224">
                thumbs up on Steam
              </Link>{' '}
              is good enough to help the site atm, but if you really want to give support you can:
              {/* Currently this site is running on someone else&apos;s server that is why you don&apos;t see ads */}
            </Typography>
          </Typography>
          <br />

          <div className={classes.item}>
            <AwardIcon className={classes.icon} color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Give the project a{' '}
              <Link
                color="secondary"
                href="https://steamcommunity.com/sharedfiles/filedetails/?id=2313234224">
                Steam Award
              </Link>{' '}
              or Like is very much appreciated.
            </Typography>
          </div>
          <br />

          <div className={classes.item}>
            <KeyIcon className={classes.icon} color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Donate{' '}
              <Link
                color="secondary"
                href="https://steamcommunity.com/tradeoffer/new/?partner=128321450&token=38BJlyuW">
                Tradable keys
              </Link>{' '}
              and will give you a donator badge on your profile in return.
            </Typography>
          </div>
          <br />

          <div className={classes.item}>
            <MoneyIcon className={classes.icon} color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Toss a coin on{' '}
              <Link color="secondary" href="/donate-coin">
                Paypal or Crypto
              </Link>{' '}
              if you want to help future cost of server and domain.
            </Typography>
          </div>
        </Container>
      </main>

      <Footer />
    </>
  )
}
