import React from 'react'
import Head from 'next/head'
import Typography from '@mui/material/Typography'
import { version } from '@/service/api'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'

export default function Version({ data }) {
  return (
    <div className="container">
      <Head>
        <meta charset="UTF-8" />
        <title>Version</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <br />
          <Typography variant="h5" component="h1">
            Version
          </Typography>
          <code>
            tag: {data.version} <br />
            hash: {data.hash} <br />
            built: {data.built} <br />
          </code>
        </Container>
      </main>

      <Footer />
    </div>
  )
}

// This gets called on every request
export async function getServerSideProps() {
  // Fetch data from external API
  // const res = await fetch(API_URL)
  // const data = await res.json()
  const data = await version()

  // Pass data to the page via props
  return { props: { data } }
}
