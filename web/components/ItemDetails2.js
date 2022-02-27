import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import MuiLink from '@mui/material/Link'
import Button from '@mui/material/Button'
import Grid from '@mui/material/Grid'
import { Box } from '@mui/system'
import { schemaOrgProduct } from '@/lib/richdata'
import { MARKET_STATUS_LIVE, MARKET_TYPE_BID } from '@/constants/market'
import { APP_NAME } from '@/constants/strings'
import { CDN_URL, marketSearch, trackItemViewURL } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import RarityTag from '@/components/RarityTag'
import MarketList from '@/components/MarketList'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import TablePagination from '@/components/TablePagination'
import ChipLink from '@/components/ChipLink'
import BuyOrderDialog from '@/components/BuyOrderDialog'
import ItemGraph from '@/components/ItemGraph'
import ItemActivity from '@/components/ItemActivity'

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
  title: {},
  media: {
    [theme.breakpoints.down('sm')]: {
      width: 300,
      height: 170,
    },
    width: 165,
    height: 110,
  },
  postItemButton: {
    [theme.breakpoints.down('sm')]: {
      margin: `8px auto !important`,
      width: '48%',
    },
    width: 165,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(1.5),
    // height: 40,
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

  const wikiLink = `https://dota2.gamepedia.com/${item.name.replace(/ +/gi, '_')}`

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
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          {/* Item Details */}
          <Grid container spacing={1.5}>
            <Grid item>
              <div style={{ background: 'rgba(0, 0, 0, 0.15)' }}>
                {item.image && (
                  <a href={wikiLink} target="_blank" rel="noreferrer noopener">
                    <ItemImage
                      className={classes.media}
                      image={item.image}
                      width={300}
                      height={170}
                      title={item.name}
                    />
                  </a>
                )}
              </div>
            </Grid>
            <Grid item>
              <Grid>
                <Typography component="h1" variant="h4">
                  {item.name}
                </Typography>
              </Grid>
              <Grid>
                <Link href={`/search?origin=${item.origin}`}>{item.origin}</Link>{' '}
                {item.rarity !== 'regular' && (
                  <>
                    &mdash;
                    <RarityTag
                      rarity={item.rarity}
                      variant="body1"
                      component={Link}
                      href={`/search?rarity=${item.rarity}`}
                    />
                  </>
                )}
              </Grid>
              <Grid>
                <Typography color="textSecondary" component="span">
                  {`Used by: `}
                </Typography>
                <Link href={`/search?hero=${item.hero}`}>{item.hero}</Link>
              </Grid>
              <Grid>
                <ChipLink label="Dota 2 Wiki" href={wikiLink} />
                &nbsp;&middot;&nbsp;
                <Typography
                  variant="body2"
                  component={MuiLink}
                  color="textPrimary"
                  href="#reserved">
                  {item.reserved_count} Reserved
                </Typography>
                &nbsp;&middot;&nbsp;
                <Typography
                  variant="body2"
                  component={MuiLink}
                  color="textPrimary"
                  href="#delivered">
                  {item.sold_count} Delivered
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          {/* Action Buttons */}
          <Box sx={{ display: 'flex', pt: 0.5, pb: 1.5 }}>
            <Button
              className={classes.postItemButton}
              sx={{ mr: 1.5 }}
              variant="outlined"
              color="secondary"
              component={Link}
              href={`/post-item?s=${item.slug}`}
              disableUnderline>
              Post this item
            </Button>
            <Button
              onClick={handleBuyOrderClick}
              className={classes.postItemButton}
              variant="outlined"
              color="bid">
              Place buy order
            </Button>
          </Box>

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
          <ItemGraph />

          {/* Lastest Activity */}
          <ItemActivity />
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
      <img src={trackItemViewURL(item.id)} alt="" />

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
