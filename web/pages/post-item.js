import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import MarketForm from '@/components/MarketForm'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function About() {
  const classes = useStyles()

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Post Item</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <MarketForm />
        </Container>
      </main>

      <Footer />
    </>
  )
}
