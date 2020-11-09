import React from 'react'
import Head from 'next/head'
import dynamic from 'next/dynamic'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'

export default function Version({ data }) {
  return (
    <div className="container">
      <Head>
        <title>version page</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <p>Your viewing version page</p>
          <p>
            tag: {data.version} <br />
            hash: {data.hash} <br />
            built: {data.built} <br />
          </p>
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
