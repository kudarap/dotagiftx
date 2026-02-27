import React, { useState } from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import Image from 'next/image'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { Grid } from '@mui/material'
import { styled } from '@mui/system'
import Link from '@/components/Link'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import SearchInput from '@/components/SearchInput'
import { APP_NAME } from '@/constants/strings'
import { heroList } from '@/service/api'

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: '#1A20278C',
  ...theme.typography.body,
  padding: theme.spacing(1),
  paddingTop: theme.spacing(1),
  textAlign: 'center',
  color: theme.palette.text.primary,
}))

const shuffleList = arr => {
  arr.sort(() => Math.random() - 0.5)
}

export default function Heroes({ heroes: allHeroes, error }) {
  const [heroes, setHeroes] = useState(allHeroes)
  const [searchTerm, setSearchTerm] = useState()
  const handleChange = term => {
    setSearchTerm(term)
    setHeroes(allHeroes.filter(v => !!v.name.match(new RegExp(term, 'gi'))))
  }

  return (
    <div className="container">
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Heroes</title>
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
          {/* <div
            style={{
              background:
                'url(https://cdn.cloudflare.steamstatic.com/steam/apps/570/library_hero.jpg?t=1724395576617) no-repeat center top',
              backgroundColor: '#2a2638ff',
              backgroundSize: 'cover',
              width: '100%',
              height: '100%',
            }}></div> */}
        </div>

        <Container style={{ position: 'relative' }}>
          {error && (
            <Typography align="center" variant="body2" color="error">
              {error}
            </Typography>
          )}

          <Typography
            sx={{ mt: -54.5, mb: 2, letterSpacing: 3, textShadow: '0 0 8px #000000b0' }}
            variant="h3"
            component="h1"
            fontWeight="bold"
            color="pimary">
            {`All Heroes (${allHeroes.length})`}
          </Typography>

          <SearchInput
            value={searchTerm}
            onChange={handleChange}
            placeholder="Search..."
            label=""
          />

          <Grid container spacing={1} sx={{ mt: 2 }}>
            {heroes.map(hero => (
              <Grid item xs={4} md={2} key={hero.name}>
                <Link href={`/search?hero=${hero.name}`} underline="none">
                  <Item>
                    <div>
                      <Image
                        src={`/assets/heroes/${hero.image}`}
                        alt={hero.name}
                        width={256 * 0.7}
                        height={144 * 0.7}
                      />
                    </div>
                    <Typography noWrap>{hero.name}</Typography>
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

Heroes.propTypes = {
  heroes: PropTypes.arrayOf(PropTypes.object),
  error: PropTypes.string,
}
Heroes.defaultProps = {
  heroes: [],
  error: null,
}

export const getStaticProps = async () => {
  const res = await heroList()
  shuffleList(res)
  return {
    props: {
      heroes: res,
    },
    revalidate: 3600, // 1 hour
  }
}
