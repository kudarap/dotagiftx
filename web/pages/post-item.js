import React from 'react'
import Head from 'next/head'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import MarketForm from '@/components/MarketForm'

export default function About() {
  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Post Item</title>
      </Head>

      <Header />

      <main sx={() =>(
        {
          [theme.breakpoints.down('md')]: {
            marginTop: theme.spacing(2),
          },
          marginTop: theme.spacing(4),
        }
      )}>
        <Container maxWidth="sm">
          <MarketForm />
        </Container>
      </main>

      <Footer />
    </>
  )
}
