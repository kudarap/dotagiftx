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

function BtcIcon(props) {
  return (
    <SvgIcon {...props}>
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
        <path
          fill="currentColor"
          d="M504 256c0 136.967-111.033 248-248 248S8 392.967 8 256 119.033 8 256 8s248 111.033 248 248zm-141.651-35.33c4.937-32.999-20.191-50.739-54.55-62.573l11.146-44.702-27.213-6.781-10.851 43.524c-7.154-1.783-14.502-3.464-21.803-5.13l10.929-43.81-27.198-6.781-11.153 44.686c-5.922-1.349-11.735-2.682-17.377-4.084l.031-.14-37.53-9.37-7.239 29.062s20.191 4.627 19.765 4.913c11.022 2.751 13.014 10.044 12.68 15.825l-12.696 50.925c.76.194 1.744.473 2.829.907-.907-.225-1.876-.473-2.876-.713l-17.796 71.338c-1.349 3.348-4.767 8.37-12.471 6.464.271.395-19.78-4.937-19.78-4.937l-13.51 31.147 35.414 8.827c6.588 1.651 13.045 3.379 19.4 5.006l-11.262 45.213 27.182 6.781 11.153-44.733a1038.209 1038.209 0 0 0 21.687 5.627l-11.115 44.523 27.213 6.781 11.262-45.128c46.404 8.781 81.299 5.239 95.986-36.727 11.836-33.79-.589-53.281-25.004-65.991 17.78-4.098 31.174-15.792 34.747-39.949zm-62.177 87.179c-8.41 33.79-65.308 15.523-83.755 10.943l14.944-59.899c18.446 4.603 77.6 13.717 68.811 48.956zm8.417-87.667c-7.673 30.736-55.031 15.12-70.393 11.292l13.548-54.327c15.363 3.828 64.836 10.973 56.845 43.035z"
        />
      </svg>
    </SvgIcon>
  )
}

function EthIcon(props) {
  return (
    <SvgIcon {...props}>
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 512">
        <path d="M311.9 260.8L160 353.6 8 260.8 160 0l151.9 260.8zM160 383.4L8 290.6 160 512l152-221.4-152 92.8z" />
      </svg>
    </SvgIcon>
  )
}

export default function Faq() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
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
