import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { isOk as checkLoggedIn, get as getLoggedInUser } from '@/service/auth'
import { catalog, CDN_URL, marketSearch, trackViewURL } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import RarityTag from '@/components/RarityTag'
import MarketList from '@/components/MarketList'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import Button from '@/components/Button'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import ContactDialog from '@/components/ContactDialog'

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

export default function ItemDetails({ item, filter, markets: initialMarkets, canonicalURL }) {
  const classes = useStyles()

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
  const metaDesc = `Buy ${item.name} from ${item.origin}${rarityText} for ${
    item.hero
  }. Price start at $${item.lowest_ask.toFixed(2)}`

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
        <meta property="og:type" content="article" />
        <meta property="og:title" content={metaTitle} />
        <meta property="og:description" content={metaDesc} />
        <meta property="og:image" content={`${CDN_URL}/${item.image}`} />
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

      <ContactDialog />
    </>
  )
}
ItemDetails.propTypes = {
  item: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
  filter: PropTypes.object,
  markets: PropTypes.object,
}
ItemDetails.defaultProps = {
  filter: {},
  markets: {
    data: [],
  },
}

const marketSearchFilter = {
  page: 1,
  status: MARKET_STATUS_LIVE,
  sort: 'price',
}

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query, req } = props

  const canonicalURL = `https://${req.headers.host}${req.url}`

  const item = await catalog(params.slug)

  const filter = { ...marketSearchFilter, item_id: item.id }
  if (query.page) {
    filter.page = Number(query.page)
  }

  let markets = {}
  let error = null
  try {
    markets = await marketSearch(filter)
  } catch (e) {
    error = e.message
  }

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
