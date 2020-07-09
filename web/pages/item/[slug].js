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
import ItemImage from '@/components/ItemImage'

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
    width: 150,
    height: 100,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(1.5),
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
          content={`Buy ${data.name} from ${
            data.origin
          } ${data.rarity.toString().toUpperCase()} for ${data.hero}. Price start at ${
            data.lowest_ask
          }`}
        />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <div className={classes.details}>
            {data.image && (
              <ItemImage
                className={classes.media}
                image={`${data.image}/300x170`}
                title={data.name}
                rarity={data.rarity}
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
