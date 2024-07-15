import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Alert from '@mui/material/Alert'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import MarketForm from '@/components/MarketForm'
import { VERIFIED_INVENTORY_VERIFIED, VERIFIED_DELIVERY_MAP_ICON } from '@/constants/verified'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function About() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Post Item</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container maxWidth="sm">
          <MarketForm />
        </Container>
      </main>

      <Footer />
    </>
  )
}
