import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import Image from 'next/image'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { Grid } from '@mui/material'
import { styled } from '@mui/system'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Link from '@/components/Link'
import { APP_NAME } from '@/constants/strings'
import { treasureList } from '@/service/api'

const LATEST_TREASURE_DROP = new Date(2025, 12, 15)

const STILL_NEW_DAYS = 30

const rarityColorMap = {
  mythical: '#8847ff',
  immortal: '#b28a33',
}

const isTreasureNew = v => {
  const releaseDate = new Date(v)
  if (!releaseDate) {
    return false
  }

  const now = new Date()
  const diff = (now - releaseDate) / (1000 * 3600 * 24)
  return diff < STILL_NEW_DAYS
}

export const isRecentTreasureNew = () => isTreasureNew(LATEST_TREASURE_DROP)

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: '#1A20278C',
  ...theme.typography.body,
  padding: theme.spacing(1),
  paddingTop: theme.spacing(2),
  textAlign: 'center',
  color: theme.palette.text.primary,
}))

export default function Treasures({ treasures, error }) {
  return (
    <div className="container">
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: All Treasures</title>
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
          }}>
          <div
            style={{
              // background:
              // 'url(https://cdn.cloudflare.steamstatic.com/steam/apps/570/library_hero.jpg?t=1724395576617) no-repeat center center',
              background:
                'url(https://cdn.steamstatic.com/apps/dota2/images/dota_react/largo/treasure_background.png) no-repeat center top',
              backgroundColor: '#2a2638ff',
              backgroundSize: 'cover',
              width: '100%',
              height: '100%',
            }}
          />
        </div>

        <Container style={{ position: 'relative' }}>
          <Typography
            sx={{ mt: -54.5, mb: 2, letterSpacing: 3, textShadow: '0 0 8px #000000b0' }}
            variant="h3"
            component="h1"
            fontWeight="bold"
            color="pimary">
            {`All Treasures (${treasures.length})`}
          </Typography>

          {error && (
            <Typography align="center" variant="body2" color="error">
              {error}
            </Typography>
          )}

          <Grid container spacing={1}>
            {treasures.map(treasure => (
              <Grid item xs={6} md={3} key={treasure.name}>
                <Link href={`/search?origin=${treasure.name}`} underline="none">
                  <Item
                    style={{
                      borderBottom: `2px solid ${rarityColorMap[treasure.rarity]}`,
                      borderTop: isTreasureNew(treasure?.release_date) ? '2px solid green' : null,
                      marginTop: isTreasureNew(treasure?.release_date) ? -2 : null,
                    }}>
                    {isTreasureNew(treasure?.release_date) && (
                      <span
                        style={{
                          position: 'absolute',
                          zIndex: 10,
                          background: 'green',
                          fontWeight: 'bolder',
                          color: 'white',
                          padding: '0 8px',
                          marginTop: -16,
                          marginLeft: -18,
                          borderBottomLeftRadius: 4,
                          borderBottomRightRadius: 4,
                        }}>
                        new
                      </span>
                    )}
                    <div>
                      <Image
                        src={`/assets/treasures/${treasure.image}`}
                        alt={treasure.name}
                        width={256}
                        height={171}
                      />
                    </div>
                    <Typography noWrap>{treasure.name}</Typography>
                  </Item>
                </Link>
              </Grid>
            ))}
          </Grid>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
Treasures.propTypes = {
  treasures: PropTypes.arrayOf(PropTypes.object),
  error: PropTypes.string,
}
Treasures.defaultProps = {
  treasures: [],
  error: null,
}

export const getStaticProps = async () => {
  const res = await treasureList()
  return {
    props: {
      treasures: res,
    },
    revalidate: 86400, // 1day
  }
}
