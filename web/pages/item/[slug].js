import React from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { item } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import RarityTag from '@/components/RarityTag'
import MarketList from '@/components/MarketList'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  details: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  media: {
    height: 100,
    marginRight: theme.spacing(1.5),
  },
}))

export default function ItemDetails({ data }) {
  const classes = useStyles()

  const router = useRouter()

  return (
    <>
      <Head>
        <title>{data.name} | Dota 2 Giftables</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <div className={classes.details}>
            <img
              className={classes.media}
              height={100}
              src="https://gamepedia.cursecdn.com/dota2_gamepedia/7/7f/Cosmetic_icon_Pipe_of_Dezun.png?version=19a51adbc336e8d2bf22b65268e4afa5"
            />
            <div>
              <Typography variant="h4">{data.name}</Typography>
              <Typography gutterBottom>
                <Typography color="textSecondary" component="span">
                  {`hero: `}
                </Typography>
                {data.hero}
                <br />

                <Typography color="textSecondary" component="span">
                  {`rarity: `}
                </Typography>
                {data.rarity === 'regular' ? (
                  data.rarity
                ) : (
                  <RarityTag rarity={data.rarity} variant="body1" component="span" />
                )}
                <br />

                <Typography color="textSecondary" component="span">
                  {`origin: `}
                </Typography>
                {data.origin}
              </Typography>
            </div>
          </div>

          <MarketList />
        </Container>
      </main>

      <Footer />
    </>
  )
}

// This gets called on every request
export async function getServerSideProps({ params }) {
  const { slug } = params
  const data = await item(slug)
  // Pass data to the page via props
  return { props: { data } }
}
