import React from 'react'
import Head from 'next/head'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { APP_NAME } from '@/constants/strings'

export default function Plus({ data }) {
  return (
    <div className="container">
      <Head>
        <title>{APP_NAME} :: Plus</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <br />
          <Typography variant="h5" component="h1">
            Dotagift Plus
          </Typography>
          <br />
          <div>List of features</div>

          <ul>
            <li>Profile Badge</li>
            <li>Refresher Orb</li>
            <li>Shopkeeper's Blessing</li>
            <li>Gem of True Sight</li>
          </ul>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
