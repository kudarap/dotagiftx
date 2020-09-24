import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { isOk as checkLoggedIn, get as getLoggedInUser } from '@/service/auth'
import { catalog, item as itemGet, CDN_URL, marketSearch, trackViewURL } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import RarityTag from '@/components/RarityTag'
import MarketList from '@/components/MarketList'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import Button from '@/components/Button'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import { APP_URL } from '@/constants/strings'
import ChipLink from '@/components/ChipLink'

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
    [theme.breakpoints.down('xs')]: {
      margin: '0 auto !important',
    },
    width: 150,
    height: 100,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(1.5),
  },
  postItemButton: {
    [theme.breakpoints.down('xs')]: {
      margin: '8px auto !important',
    },
    width: 150,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(1.5),
  },
}))

export default function ItemDetails({
  item,
  error: initialError,
  filter,
  markets: initialMarkets,
  canonicalURL,
}) {
  const classes = useStyles()

  if (initialError) {
    return (
      <>
        <Header />

        <main className={classes.main}>
          <Container>
            <Typography variant="h5" component="h1" gutterBottom align="center">
              Item Error
            </Typography>
            <Typography color="textSecondary" align="center">
              {initialError}
            </Typography>
          </Container>
        </main>

        <Footer />
      </>
    )
  }

  const [page, setPage] = React.useState(filter.page)
  const [markets, setMarkets] = React.useState(initialMarkets)
  const [error, setError] = React.useState(null)

  // Handle market request on page change.
  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await marketSearch({ ...filter, page })
        setMarkets(res)
      } catch (e) {
        setError(e.message)
      }
    })()
  }, [page])

  const handlePageChange = (e, p) => {
    setPage(p)
  }

  const linkProps = { href: '/item/[slug]', as: `/item/${item.slug}` }

  const metaTitle = `DotagiftX :: Listings for ${item.name}`
  const rarityText = item.rarity === 'regular' ? '' : ` — ${item.rarity.toString().toUpperCase()}`
  let metaDesc = `Buy ${item.name} from ${item.origin}${rarityText} item for ${item.hero}.`
  const schemaOrgProd = {
    '@context': 'https://schema.org',
    '@type': 'Product',
    productID: item.id,
    name: item.name,
    image: `${CDN_URL}/${item.image}`,
    description: metaDesc,
    offer: {
      '@type': 'Offer',
      priceCurrency: 'USD',
      url: canonicalURL,
    },
  }
  if (item.lowest_ask) {
    const startingPrice = item.lowest_ask.toFixed(2)
    metaDesc += ` Price starting at $${startingPrice}`
    schemaOrgProd.offers = {
      availability: 'https://schema.org/InStock',
      price: startingPrice,
    }
  } else {
    schemaOrgProd.offers = {
      availability: 'https://schema.org/OutOfStock',
      price: '0',
    }
  }

  const isLoggedIn = checkLoggedIn()
  let currentUserID = null
  if (isLoggedIn) {
    currentUserID = getLoggedInUser().user_id
  }

  return (
    <>
      <Head>
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={canonicalURL} />

        {/* Twitter Card */}
        <meta name="twitter:card" content="summary" />
        <meta name="twitter:title" content={metaTitle} />
        <meta name="twitter:description" content={metaDesc} />
        <meta name="twitter:image" content={`${CDN_URL}/${item.image}`} />
        <meta name="twitter:site" content="@DotagiftX" />
        {/* OpenGraph */}
        <meta property="og:url" content={canonicalURL} />
        <meta property="og:type" content="website" />
        <meta property="og:title" content={metaTitle} />
        <meta property="og:description" content={metaDesc} />
        <meta property="og:image" content={`${CDN_URL}/${item.image}`} />
        {/* Rich Results */}
        <script type="application/ld+json">{JSON.stringify(schemaOrgProd)}</script>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <div className={classes.details}>
            {item.image && (
              <div>
                <ItemImage
                  className={classes.media}
                  image={`/300x170/${item.image}`}
                  title={item.name}
                  rarity={item.rarity}
                />
                {isLoggedIn && (
                  <Button
                    className={classes.postItemButton}
                    variant="outlined"
                    color="secondary"
                    size="small"
                    component={Link}
                    href={`/post-item?s=${item.slug}`}
                    fullWidth>
                    Post this Item
                  </Button>
                )}
              </div>
            )}

            <Typography component="h1">
              <Typography component="p" variant="h4">
                {item.name}
              </Typography>
              <Typography gutterBottom>
                <Link href={`/search?q=${item.origin}`}>{item.origin}</Link>{' '}
                {item.rarity !== 'regular' && (
                  <>
                    &mdash;
                    <RarityTag rarity={item.rarity} variant="body1" component="span" />
                  </>
                )}
                <br />
                <Typography color="textSecondary" component="span">
                  {`Used by: `}
                </Typography>
                <Link href={`/search?q=${item.hero}`}>{item.hero}</Link>
                <br />
                <ChipLink
                  label="Check Dota 2 Wiki"
                  href={`https://dota2.gamepedia.com/${item.name}`}
                />
              </Typography>
            </Typography>
          </div>

          <MarketList data={markets} currentUserID={currentUserID} error={error} />
          {!error && (
            <TablePaginationRouter
              linkProps={linkProps}
              style={{ textAlign: 'right' }}
              count={markets.total_count}
              page={page}
              onChangePage={handlePageChange}
            />
          )}
        </Container>

        <img src={trackViewURL(item.id)} alt="" />
      </main>

      <Footer />
    </>
  )
}
ItemDetails.propTypes = {
  item: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
  filter: PropTypes.object,
  markets: PropTypes.object,
  error: PropTypes.string,
}
ItemDetails.defaultProps = {
  filter: {},
  markets: {
    data: [],
  },
  error: null,
}

const marketSearchFilter = {
  page: 1,
  status: MARKET_STATUS_LIVE,
  sort: 'price',
}

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query } = props

  // Handles invalid item slug
  let item = {}
  try {
    item = await itemGet(params.slug)
  } catch (e) {
    return {
      props: {
        item,
        error: e.message,
        filter: {},
        markets: {},
      },
    }
  }

  // Handles no market entry on item
  try {
    item = await catalog(params.slug)
  } catch (e) {
    console.log(`catalog get error: ${e.message}`)
  }
  if (!item.id) {
    return {
      props: {
        item,
        filter: {},
      },
    }
  }

  const filter = { ...marketSearchFilter, item_id: item.id }
  if (query.page) {
    filter.page = Number(query.page)
  }

  let markets = {}
  let error = null
  try {
    markets = await marketSearch(filter)
  } catch (e) {
    console.log(`market search error: ${e.message}`)
    error = e.message
  }

  const canonicalURL = `${APP_URL}/item/${params.slug}`

  return {
    props: {
      item,
      canonicalURL,
      filter,
      markets,
      error,
    },
  }
}
