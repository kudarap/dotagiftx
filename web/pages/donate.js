import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import Footer from '@/components/Footer'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(6),
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
            Share your Tango!
          </Typography>
          <br />

          <Typography color="textSecondary">
            <img src="/assets/tango.png" style={{ float: 'right' }} alt="tango" />
            {/* I had this idea once to create a community market for Dota 2 giftables items and now */}
            {/* here we are, although the idea is free but the time to develop and server are not and */}
            {/* running on someone else&apos;s server 🤫. */}
            If this project helped you somehow and want to support it, you can toss a tango to your
            developer.
          </Typography>
          <br />

          <div>
            <form action="https://www.paypal.com/donate" method="post" target="_top">
              <input type="hidden" name="cmd" value="_donations" />
              <input type="hidden" name="business" value="LBY7VY8PQ9D3Y" />
              <input type="hidden" name="currency_code" value="USD" />
              <input
                type="image"
                src="/assets/tango.png"
                border="0"
                name="submit"
                title="PayPal - The safer, easier way to pay online!"
                alt="Donate with PayPal button"
              />
              <img
                alt=""
                border="0"
                src="https://www.paypal.com/en_PH/i/scr/pixel.gif"
                width="1"
                height="1"
              />
            </form>
          </div>

          <Typography variant="h6" component="h2">
            How can I report scammers?
          </Typography>
          <Typography color="textSecondary">
            You can use <Link href="https://steamrep.com/">SteamRep</Link> or inquire on{' '}
            <Link href="https://www.reddit.com/r/Dota2Trade/">r/Dota2Trade</Link>.
          </Typography>
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
