import React from 'react'
import Head from 'next/head'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { Grid } from '@mui/material'
import { styled } from '@mui/system'
import Image from 'next/image'

import placeholder from '../public/assets/treasure.png'

const treasures = [
  {
    name: "Aghanim's 2021 Ageless Heirlooms",
    image: 'placeholder.png',
    rarity: 'mythical',
    items: 10,
  },
  {
    name: "Aghanim's 2021 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'mythical',
    items: 17,
  },
  {
    name: "Aghanim's 2021 Continuum Collection",
    image: 'placeholder.png',
    rarity: 'mythical',
    items: 7,
  },
  {
    name: "Aghanim's 2021 Immortal Treasure",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 9,
  },
  {
    name: 'Immortal Treasure I 2020',
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 10,
  },
  {
    name: 'Immortal Treasure II 2020',
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 10,
  },
  {
    name: 'Immortal Treasure III 2020',
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 8,
  },
  {
    name: "Nemestice 2021 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 15,
  },
  {
    name: 'Nemestice 2021 Immortal Treasure',
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 9,
  },
  {
    name: 'Nemestice 2021 Themed Treasure',
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 11,
  },
  {
    name: "The International 2015 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 11,
  },
  {
    name: "The International 2016 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 15,
  },
  {
    name: "The International 2017 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 22,
  },
  {
    name: "The International 2018 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 17,
  },
  {
    name: "The International 2018 Collector's Cache II",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 14,
  },
  {
    name: "The International 2019 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 18,
  },
  {
    name: "The International 2019 Collector's Cache II",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 16,
  },
  {
    name: "The International 2020 Collector's Cache",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 18,
  },
  {
    name: "The International 2020 Collector's Cache II",
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 17,
  },
  {
    name: 'Treasure of the Cryptic Beacon',
    image: 'placeholder.png',
    rarity: 'immortal',
    items: 6,
  },
]

const rarityColorMap = {
  mythical: '#8847ff',
  immortal: '#b28a33',
}

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: '#1A2027CC',
  ...theme.typography.body,
  padding: theme.spacing(1),
  textAlign: 'center',
  color: theme.palette.text.secondary,
}))

const pageTitle = 'All Giftable Treasures'

export default function Treasures({ data }) {
  return (
    <div className="container">
      <Head>
        <title>{pageTitle}</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <div
          style={{
            width: '100%',
            height: 500,
            maskImage: 'linear-gradient(to top, transparent 25%, black 100%)',
            position: 'relative',
            zIndex: 0,
          }}>
          <div
            style={{
              background: 'url(/assets/treasure-banner.png) no-repeat top center',
              width: '100%',
              height: '100%',
            }}></div>
        </div>

        <Container style={{ position: 'relative' }}>
          <Typography variant="h5" component="h1" sx={{ mt: -35, mb: 4 }}>
            {pageTitle}
          </Typography>

          <Grid container spacing={1}>
            {treasures.map(treasure => {
              return (
                <Grid item xs={6} md={4}>
                  <Item style={{ borderBottom: `2px solid ${rarityColorMap[treasure.rarity]}` }}>
                    <div>
                      <Image src={placeholder} alt={treasure.name} layout="responsive" />
                    </div>
                    {treasure.name}
                  </Item>
                </Grid>
              )
            })}
          </Grid>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
