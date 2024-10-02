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
import Link from '@/components/Link'
import { APP_NAME } from '@/constants/strings'

const treasures = [
  {
    name: "Crownfall 2024 Collector's Cache II",
    image: 'crownfall_2024_collect_s_cache_ii.png',
    rarity: 'mythical',
    items: 16,
  },
  {
    name: "Crownfall 2024 Collector's Cache",
    image: 'crownfall_2024_collect_s_cache.png',
    rarity: 'mythical',
    items: 16,
  },
  {
    name: 'Crownfall Treasure I',
    image: 'crownfall_treasure_i.png',
    rarity: 'mythical',
    items: 12,
  },
  {
    name: 'Crownfall Treasure II',
    image: 'crownfall_treasure_ii.png',
    rarity: 'mythical',
    items: 12,
  },
  {
    name: 'Crownfall Treasure III',
    image: 'crownfall_treasure_iii.png',
    rarity: 'mythical',
    items: 11,
  },
  {
    name: 'Dead Reckoning Chest',
    image: 'dead_reckoning_chest.png',
    rarity: 'mythical',
    items: 14,
  },
  {
    name: "August 2023 Collector's Cache",
    image: 'august_2023_collector_s_cache.png',
    rarity: 'mythical',
    items: 16,
  },
  {
    name: "Diretide 2022 Collector's Cache II",
    image: 'diretide_2022_collector_s_cache_ii.png',
    rarity: 'immortal',
    items: 17,
  },
  {
    name: "Diretide 2022 Collector's Cache",
    image: 'diretide_2022_collector_s_cache.png',
    rarity: 'immortal',
    items: 18,
  },
  {
    name: 'Immortal Treasure I 2022',
    image: 'immortal_treasure_i_2022.png',
    rarity: 'immortal',
    items: 9,
  },
  {
    name: 'Immortal Treasure II 2022',
    image: 'immortal_treasure_ii_2022.png',
    rarity: 'immortal',
    items: 9,
  },
  {
    name: 'The Battle Pass Collection 2022',
    image: 'the_battle_pass_collection_2022.png',
    rarity: 'immortal',
    items: 8,
  },
  {
    name: 'Ageless Heirlooms 2022',
    image: 'ageless_heirlooms_2022.png',
    rarity: 'immortal',
    items: 10,
  },
  {
    name: "Aghanim's 2021 Collector's Cache",
    image: 'aghanim_s_2021_collector_s_cache.webp',
    rarity: 'mythical',
    items: 17,
  },
  {
    name: "Aghanim's 2021 Ageless Heirlooms",
    image: 'aghanim_s_2021_ageless_heirlooms.webp',
    rarity: 'mythical',
    items: 10,
  },
  {
    name: "Aghanim's 2021 Continuum Collection",
    image: 'aghanim_s_2021_continuum_collection.webp',
    rarity: 'mythical',
    items: 7,
  },
  {
    name: "Aghanim's 2021 Immortal Treasure",
    image: 'aghanim_s_2021_immortal_treasure.webp',
    rarity: 'immortal',
    items: 9,
  },
  {
    name: "Nemestice 2021 Collector's Cache",
    image: 'nemestice_2021_collector_s_cache.webp',
    rarity: 'mythical',
    items: 15,
  },
  {
    name: 'Nemestice 2021 Immortal Treasure',
    image: 'nemestice_2021_immortal_treasure.webp',
    rarity: 'mythical',
    items: 9,
  },
  {
    name: 'Nemestice 2021 Themed Treasure',
    image: 'nemestice_2021_themed_treasure.webp',
    rarity: 'mythical',
    items: 11,
  },
  {
    name: 'Immortal Treasure I 2020',
    image: 'immortal_treasure_i_2020.webp',
    rarity: 'immortal',
    items: 10,
  },
  {
    name: 'Immortal Treasure II 2020',
    image: 'immortal_treasure_ii_2020.webp',
    rarity: 'immortal',
    items: 10,
  },
  {
    name: 'Immortal Treasure III 2020',
    image: 'immortal_treasure_iii_2020.webp',
    rarity: 'immortal',
    items: 8,
  },
  {
    name: "The International 2020 Collector's Cache",
    image: 'the_international_2020_collector_s_cache.webp',
    rarity: 'mythical',
    items: 18,
  },
  {
    name: "The International 2020 Collector's Cache II",
    image: 'the_international_2020_collector_s_cache_ii.webp',
    rarity: 'mythical',
    items: 17,
  },
  {
    name: "The International 2019 Collector's Cache",
    image: 'the_international_2019_collector_s_cache.webp',
    rarity: 'mythical',
    items: 18,
  },
  {
    name: "The International 2019 Collector's Cache II",
    image: 'the_international_2019_collector_s_cache_ii.webp',
    rarity: 'mythical',
    items: 16,
  },
  {
    name: "The International 2018 Collector's Cache",
    image: 'the_international_2018_collector_s_cache.webp',
    rarity: 'mythical',
    items: 17,
  },
  {
    name: "The International 2018 Collector's Cache II",
    image: 'the_international_2018_collector_s_cache_ii.webp',
    rarity: 'mythical',
    items: 14,
  },
  {
    name: "The International 2017 Collector's Cache",
    image: 'the_international_2017_collector_s_cache.webp',
    rarity: 'mythical',
    items: 22,
  },
  {
    name: "The International 2016 Collector's Cache",
    image: 'the_international_2016_collector_s_cache.webp',
    rarity: 'mythical',
    items: 15,
  },
  {
    name: "The International 2015 Collector's Cache",
    image: 'the_international_2015_collector_s_cache.webp',
    rarity: 'mythical',
    items: 11,
  },
  {
    name: 'Treasure of the Cryptic Beacon',
    image: 'treasure_of_the_cryptic_beacon.webp',
    rarity: 'mythical',
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

export default function Treasures() {
  return (
    <div className="container">
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: All Giftable Treasures</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <div
          style={{
            width: '100%',
            height: 500,
            maskImage: 'linear-gradient(to top, transparent 25%, black 100%)',
            WebkitMaskImage: 'linear-gradient(to top, transparent 0%, black 90%)',
            position: 'relative',
            zIndex: 0,
            top: -7,
          }}>
          <div
            style={{
              background: 'url(/assets/venge_arcana.png) no-repeat bottom center',
              backgroundPosition: 'center',
              width: '100%',
              height: '100%',
            }}></div>
        </div>

        <Container style={{ position: 'relative' }}>
          <Typography variant="h4" component="h1" sx={{ mt: -35, mb: 4 }}>
            All Giftable Treasures
          </Typography>

          <Grid container spacing={1}>
            {treasures.map(treasure => {
              return (
                <Grid item xs={6} md={3}>
                  <Link href={`/search?origin=${treasure.name}`} underline="none">
                    <Item style={{ borderBottom: `2px solid ${rarityColorMap[treasure.rarity]}` }}>
                      <div>
                        <Image
                          src={'/assets/treasures/' + treasure.image}
                          alt={treasure.name}
                          width={256}
                          height={171}
                        />
                      </div>
                      <Typography noWrap>{treasure.name}</Typography>
                    </Item>
                  </Link>
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
