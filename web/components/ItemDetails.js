import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import useInView from 'react-cool-inview'
import Typography from '@mui/material/Typography'
import Button from '@mui/material/Button'
import Grid from '@mui/material/Grid'
import { schemaOrgProduct } from '@/lib/richdata'
import { MARKET_STATUS_LIVE, MARKET_TYPE_BID } from '@/constants/market'
import { APP_NAME } from '@/constants/strings'
import { CDN_URL, marketSearch, trackItemViewURL } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import MarketList from '@/components/MarketList'
import Link from '@/components/Link'
import TablePagination from '@/components/TablePagination'
import BuyOrderDialog from '@/components/BuyOrderDialog'
import ItemViewCard from '@/components/ItemViewCard'
import dynamic from 'next/dynamic'

const ItemGraph = dynamic(() => import('@/components/ItemGraph'))

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  details: {
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  postItemButton: {
    [theme.breakpoints.down('sm')]: {
      width: '50%',
    },
    width: 172,
  },
}))

const marketBuyOrderFilter = {
  type: MARKET_TYPE_BID,
  status: MARKET_STATUS_LIVE,
  sort: 'highest',
  nocache: true,
}

const DEFAULT_SORT = 'price'
const OFFERS_PARAM_KEY = 'offers'
const BUYORDERS_PARAM_KEY = 'buyorders'

export default function ItemDetails({
  item,
  error: initialError,
  filter: initialFilter,
  marketType,
  sortParam,
  initialAsks,
  initialBids,
  canonicalURL,
}) {
  const { classes } = useStyles()

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

  const [offers, setOffers] = React.useState(initialAsks)
  const [orders, setOrders] = React.useState(initialBids)
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(null)
  const [openBuyOrderDialog, setOpenBuyOrderDialog] = React.useState(false)
  const [tabIndex, setTabIndex] = React.useState(0)

  const router = useRouter()

  // Set active tab on load
  React.useEffect(() => {
    switch (marketType) {
      case OFFERS_PARAM_KEY:
        setTabIndex(0)
        break
      case BUYORDERS_PARAM_KEY:
        setTabIndex(1)
        break
      default:
        setTabIndex(0)
    }
  }, [marketType])

  // Handle offers data on load. when its available display immediately.
  React.useEffect(() => {
    setOffers(initialAsks)
  }, [initialAsks])

  // Handle filter changes
  const [sort, setSort] = React.useState(sortParam || DEFAULT_SORT)
  const [page, setPage] = React.useState(initialFilter.page)
  const handleTabChange = idx => {
    setTabIndex(idx)
    setSort('price')
    setPage(1)

    // process offers
    if (idx === 0) {
      router.push(`/${item.slug}/${OFFERS_PARAM_KEY}`, null, { shallow: true })
      return
    }

    // process buy orders
    router.push(`/${item.slug}/${BUYORDERS_PARAM_KEY}`, null, { shallow: true })
  }
  const handleSortChange = sortValue => {
    setSort(sortValue)
    setPage(1)

    // process offers
    if (tabIndex === 0) {
      router.push(`/${item.slug}/${OFFERS_PARAM_KEY}/${sortValue}`, null, { shallow: true })
      return
    }

    // process buy orders
    router.push(`/${item.slug}/${BUYORDERS_PARAM_KEY}/${sortValue}`, null, { shallow: true })
  }
  const handlePageChange = (e, pageValue) => {
    setPage(pageValue)
    // scroll to top when page change
    window.scrollTo(0, 0)
    router.push({ query: { page: pageValue } }, null, { shallow: true })
  }

  const getOffers = async (sortValue, pageValue) => {
    setLoading('ask')
    try {
      const res = await marketSearch({
        ...initialFilter,
        sort: sortValue === 'price' ? 'lowest' : sortValue,
        page: pageValue,
        index: 'item_id',
      })
      setOffers(res)
    } catch (e) {
      setError(e.message)
    }
    setLoading(null)
  }
  const getBuyOrders = async sortValue => {
    setLoading('bid')
    try {
      const res = await marketSearch({
        ...marketBuyOrderFilter,
        sort: sortValue === 'price' ? 'highest' : sortValue,
        item_id: item.id,
        index: 'item_id',
      })
      res.loaded = true
      setOrders(res)
    } catch (e) {
      setError(e.message)
    }
    setLoading(null)
  }
  const handleBuyOrderClick = () => {
    setOpenBuyOrderDialog(true)
  }
  const handleBuyerChange = () => {
    getBuyOrders(sort)
  }

  // Handle initial buy orders on page load.
  React.useEffect(() => {
    getBuyOrders(sort)
  }, [])

  // Handles update offers and buy orders on filter change
  React.useEffect(() => {
    // process offers
    if (tabIndex === 0) {
      // check initial props is same and skip the fetch
      if (sort === sortParam && page === initialFilter.page) {
        setOffers(initialAsks)
        return
      }

      getOffers(sort, page)
      return
    }

    // process buy orders
    getBuyOrders(sort)
  }, [tabIndex, sort, page])

  const metaTitle = `${APP_NAME} :: Listings for ${item.name}`
  const rarityText = item.rarity === 'regular' ? '' : ` — ${item.rarity.toString().toUpperCase()}`
  let metaDesc = `Buy ${item.name} from ${item.origin}${rarityText} item for ${item.hero}.`
  const jsonLD = schemaOrgProduct(canonicalURL, item, { description: metaDesc })
  if (item.lowest_ask) {
    const startingPrice = item.lowest_ask.toFixed(2)
    metaDesc += ` Price starting at $${startingPrice}`
  }

  const { observe, inView } = useInView({
    onEnter: ({ unobserve }) => unobserve(), // only run once
  })

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={canonicalURL} />

        {/* Twitter Card */}
        <meta name="twitter:card" content="summary" />
        <meta name="twitter:title" content={metaTitle} />
        <meta name="twitter:description" content={metaDesc} />
        <meta name="twitter:image" content={`${CDN_URL}/${item.image}`} />
        <meta name="twitter:site" content={`${APP_NAME}`} />
        {/* OpenGraph */}
        <meta property="og:url" content={canonicalURL} />
        <meta property="og:type" content="website" />
        <meta property="og:title" content={metaTitle} />
        <meta property="og:description" content={metaDesc} />
        <meta property="og:image" content={`${CDN_URL}/${item.image}`} />
        {/* Rich Results */}
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLD) }}
        />

        {/* Preload the LCP image with a high fetchpriority so it starts loading with the stylesheet. */}
        <link
          rel="preload"
          fetchpriority="high"
          as="image"
          href={`${CDN_URL}/${item.image}`}
          type="image/png"
        />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          {/* Item Details */}
          <ItemViewCard item={item} />

          {/* Action Buttons */}
          <Grid container alignItems="center" spacing={1} sx={{ mb: 3 }}>
            <Grid item className={classes.postItemButton}>
              <Button
                fullWidth
                variant="outlined"
                color="secondary"
                component={Link}
                href={`/post-item?s=${item.slug}`}
                disableUnderline>
                Post this item
              </Button>
            </Grid>
            <Grid item className={classes.postItemButton}>
              <Button fullWidth onClick={handleBuyOrderClick} variant="outlined" color="bid">
                Place buy order
              </Button>
            </Grid>
          </Grid>

          {/* Listings */}
          <MarketList
            offers={offers}
            buyOrders={orders}
            error={error}
            loading={loading}
            onSortChange={handleSortChange}
            tabIndex={tabIndex}
            onTabChange={handleTabChange}
            sort={sort}
            pagination={
              !error && (
                <TablePagination
                  onPageChange={handlePageChange}
                  style={{ textAlign: 'right' }}
                  count={offers.total_count || 0}
                  page={page}
                />
              )
            }
          />

          {/* History */}
          <div ref={observe}>{inView && <ItemGraph itemId={item.id} itemName={item.name} />}</div>
        </Container>
      </main>

      <BuyOrderDialog
        catalog={item}
        open={openBuyOrderDialog}
        onClose={() => {
          setOpenBuyOrderDialog(false)
        }}
        onChange={handleBuyerChange}
      />
      <img src={trackItemViewURL(item.id)} height={1} width={1} alt="" />

      <Footer />
    </>
  )
}
ItemDetails.propTypes = {
  item: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
  marketType: PropTypes.string.isRequired,
  filter: PropTypes.object,
  sortParam: PropTypes.string,
  initialAsks: PropTypes.object,
  initialBids: PropTypes.object,
  error: PropTypes.string,
}
ItemDetails.defaultProps = {
  filter: {},
  sortParam: DEFAULT_SORT,
  initialAsks: {
    data: [],
  },
  initialBids: {
    data: [],
  },
  error: null,
}
