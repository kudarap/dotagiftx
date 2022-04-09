import React from 'react'
import Head from 'next/head'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'

export default function Version({ data }) {
  return (
    <div className="container">
      <Head>
        <title>Treasury</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <br />
          <Typography variant="h5" component="h1">
            Treasury
          </Typography>
          <br />
          <div>List of Treaures</div>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
