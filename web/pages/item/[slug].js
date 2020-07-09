import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { CDN_URL, catalog, trackViewURL } from '@/service/api'
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

  return (
    <>
      <Head>
        <title>
          Dota 2 Giftables :: Listings for {data.name} starts at ${data.lowest_ask}
        </title>
        <meta
          name="description"
          content={`Buy ${data.name} ${data.rarity}. Price at ${data.lowest_ask}`}
        />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <div className={classes.details}>
            {data.image && (
              <img
                className={classes.media}
                height={100}
                alt={data.name}
                src={`${CDN_URL + data.image}/300x170`}
              />
            )}
            <Typography component="h1">
              <Typography variant="h4">{data.name}</Typography>
              <Typography gutterBottom>
                {data.origin}{' '}
                {data.rarity !== 'regular' && (
                  <>
                    &mdash;
                    <RarityTag rarity={data.rarity} variant="body1" component="span" />
                  </>
                )}
                <br />
                <Typography color="textSecondary" component="span">
                  {`Used by: `}
                </Typography>
                {data.hero}
              </Typography>
            </Typography>
          </div>

          <MarketList itemID={data.id} />
        </Container>

        <img src={trackViewURL(data.id)} alt="" />
      </main>

      <Footer />
    </>
  )
}
ItemDetails.propTypes = {
  data: PropTypes.object,
}
ItemDetails.defaultProps = {
  data: {},
}

// This gets called on every request
export async function getServerSideProps({ params }) {
  const { slug } = params
  const data = await catalog(slug)
  // Pass data to the page via props
  return { props: { data } }
}
