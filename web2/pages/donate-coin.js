import React from 'react'
import Head from 'next/head'
import { makeStyles, withStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import SvgIcon from '@mui/material/SvgIcon'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Button from '@/components/Button'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  grid: {
    margin: theme.spacing(3, 0, 6),
  },
  crypto: {
    margin: theme.spacing(2, 8, 0, 0),
    '& img': {
      margin: theme.spacing(2, 0, 0, 2),
    },
  },
}))

const PaypalButton = withStyles(Button, theme => ({
  root: {
    marginTop: theme.spacing(2),
    width: 300,
    color: theme.palette.getContrastText('#0070ba'),
    backgroundColor: '#0070ba',
    '&:hover': {
      color: theme.palette.getContrastText('#ffc439'),
      backgroundColor: '#ffc439',
    },
  },
}))

function PaypalIcon(props) {
  return (
    <SvgIcon {...props}>
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512">
        <path
          fill="currentColor"
          d="M111.4 295.9c-3.5 19.2-17.4 108.7-21.5 134-.3 1.8-1 2.5-3 2.5H12.3c-7.6 0-13.1-6.6-12.1-13.9L58.8 46.6c1.5-9.6 10.1-16.9 20-16.9 152.3 0 165.1-3.7 204 11.4 60.1 23.3 65.6 79.5 44 140.3-21.5 62.6-72.5 89.5-140.1 90.3-43.4.7-69.5-7-75.3 24.2zM357.1 152c-1.8-1.3-2.5-1.8-3 1.3-2 11.4-5.1 22.5-8.8 33.6-39.9 113.8-150.5 103.9-204.5 103.9-6.1 0-10.1 3.3-10.9 9.4-22.6 140.4-27.1 169.7-27.1 169.7-1 7.1 3.5 12.9 10.6 12.9h63.5c8.6 0 15.7-6.3 17.4-14.9.7-5.4-1.1 6.1 14.4-91.3 4.6-22 14.3-19.7 29.3-19.7 71 0 126.4-28.8 142.9-112.3 6.5-34.8 4.6-71.4-23.8-92.6z"
        />
      </svg>
    </SvgIcon>
  )
}

export default function Faq() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Donate</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Toss a coin
          </Typography>
          <br />

          <Typography color="textSecondary">
            <img src="/assets/midas.png" style={{ float: 'right' }} alt="tango" />I had this idea
            once to create a community market for Dota 2 Giftables items and now here we are,
            although the idea is free(had fun writing it) and server isn&apos;t. BUT! thanks to
            someone else&apos;s server running this website. <br />
            <br />
            If this project helped you sold your items or struck a good deal and really want to
            support it, you can toss a coin to help cost for the server.
          </Typography>
          <br />

          <PaypalButton
            startIcon={<PaypalIcon />}
            size="large"
            target="_blank"
            rel="noreferrer noopener"
            href="https://www.paypal.com/donate?hosted_button_id=QHWKBTN4VGDR6">
            Donate with PayPal
          </PaypalButton>
          <Typography color="textSecondary" style={{ marginTop: 6 }}>
            <strike>
              Please don&apos;t forget to put your profile link on the notes so I can award the
              badge
            </strike>
          </Typography>

          {/* <Grid container className={classes.grid} alignContent="center">
            <Grid item className={classes.crypto}>
              <Typography variant="body2">
                <BtcIcon fontSize="inherit" /> <strong>Bitcoin(BTC)</strong>
                <br />
                3QH7ofHgxoUNu2X4v29KA9inoP6ELFchVB
              </Typography>
              <img src="/assets/btc_qr.png" alt="Bitcoin QR" />
            </Grid>
            <Grid item className={classes.crypto}>
              <Typography variant="body2">
                <EthIcon fontSize="inherit" /> <strong>Etherium(ETH)</strong>
                <br />
                0x928cE680130328eb9578Fd22413507B99fe1135C
              </Typography>
              <img src="/assets/eth_qr.png" alt="Bitcoin QR" />
            </Grid>
          </Grid> */}
        </Container>
      </main>

      <Footer />
    </>
  )
}
