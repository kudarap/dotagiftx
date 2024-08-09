import React, { useEffect } from 'react'
import Head from 'next/head'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'

const GIVEAWAY_LINK = 'https://gamesandanimes.site/giveaway-dota-2/'

export default function Giveaway() {
  useEffect(() => {
    if (window !== undefined) {
      return
    }
    window.location = GIVEAWAY_LINK
  }, [])

  return (
    <div className="container">
      <Head>
        <meta charset="UTF-8" />
        <title>Giveaway</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <br />
          <Typography variant="h5" component="h1">
            Giveaway
          </Typography>
          <Typography sx={{ minHeight: '50vh' }}>Redirecting...</Typography>
        </Container>
      </main>

      <Footer />
    </div>
  )
}

// This gets called on every request
export async function getServerSideProps() {
  return {
    redirect: {
      permanent: false,
      destination: GIVEAWAY_LINK,
    },
  }
}
