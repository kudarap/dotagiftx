import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE, MARKET_TYPE_BID } from '@/constants/market'
import { itemRarityColorMap } from '@/constants/palette'
import { APP_NAME } from '@/constants/strings'
import { CDN_URL, marketSearch, trackViewURL } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import RarityTag from '@/components/RarityTag'
import MarketList from '@/components/MarketList'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import Button from '@/components/Button'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import ChipLink from '@/components/ChipLink'
import AppContext from '@/components/AppContext'
import BidButton from '@/components/BidButton'
import BuyOrderDialog from '@/components/BuyOrderDialog'

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
  title: {},
  media: {
    [theme.breakpoints.down('xs')]: {
      margin: '8px auto 8px !important',
      width: 300,
      height: 170,
    },
    width: 165,
    height: 110,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(1.5),
  },
  postItemButton: {
    [theme.breakpoints.down('xs')]: {
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
  sort: 'price:desc',
}

export default function ItemDetails({
  item,
  error: initialError,
  filter,
  markets: initialMarkets,
  canonicalURL,
}) {
  const classes = useStyles()

  const { isMobile, isLoggedIn } = useContext(AppContext)

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

  const [markets, setMarkets] = React.useState(initialMarkets)
  const [buyOrders, setBuyOrders] = React.useState(initialMarkets)
  const [error, setError] = React.useState(null)
  const [openBuyOrderDialog, setOpenBuyOrderDialog] = React.useState(false)

  // Handle market request on page change.
  marketBuyOrderFilter.item_id = item.id
  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await marketSearch(filter)
        setMarkets(res)
        const res2 = await marketSearch(marketBuyOrderFilter)
        setBuyOrders(res2)
      } catch (e) {
        setError(e.message)
      }
    })()
  }, [filter])

  const handleBuyOrderClick = () => {
    setOpenBuyOrderDialog(true)
  }

  const metaTitle = `${APP_NAME} :: Listings for ${item.name}`
  const rarityText = item.rarity === 'regular' ? '' : ` — ${item.rarity.toString().toUpperCase()}`
  let metaDesc = `Buy ${item.name} from ${item.origin}${rarityText} item for ${item.hero}.`
  const schemaOrgProd = {
    '@context': 'https://schema.org',
    '@type': 'Product',
    productID: item.id,
    name: item.name,
    image: `${CDN_URL}/${item.image}`,
    description: metaDesc,
    offers: {
      '@type': 'Offer',
      priceCurrency: 'USD',
      url: canonicalURL,
    },
  }
  if (item.lowest_ask) {
    const startingPrice = item.lowest_ask.toFixed(2)
    metaDesc += ` Price starting at $${startingPrice}`
    schemaOrgProd.offers.availability = 'https://schema.org/InStock'
    schemaOrgProd.offers.price = startingPrice
  } else {
    schemaOrgProd.offers.availability = 'https://schema.org/OutOfStock'
    schemaOrgProd.offers.price = '0'
  }

  const wikiLink = `https://dota2.gamepedia.com/${item.name.replace(/ +/gi, '_')}`
  const linkProps = { href: `/${item.slug}` }

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
          dangerouslySetInnerHTML={{ __html: JSON.stringify(schemaOrgProd) }}
        />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          {!isMobile ? (
            <div className={classes.details}>
              {item.image && (
                <div>
                  <a href={wikiLink} target="_blank" rel="noreferrer noopener">
                    <ItemImage
                      className={classes.media}
                      image={item.image}
                      width={165}
                      height={110}
                      title={item.name}
                      rarity={item.rarity}
                    />
                  </a>
                  <Button
                    className={classes.postItemButton}
                    variant="outlined"
                    color="secondary"
                    component={Link}
                    href={`/post-item?s=${item.slug}`}
                    disableUnderline
                    fullWidth>
                    Post this Item
                  </Button>
                </div>
              )}

              <Typography component="h1">
                <Typography component="p" variant="h4">
                  {item.name}
                </Typography>
                <Typography gutterBottom>
                  <Link href={`/search?origin=${item.origin}`} color="textSecondary">
                    {item.origin}
                  </Link>{' '}
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
                  <br />
                  <Typography color="textSecondary" component="span">
                    {`Used by: `}
                  </Typography>
                  <Link color="textSecondary" href={`/search?hero=${item.hero}`}>
                    {item.hero}
                  </Link>
                  <br />
                  <ChipLink label="Dota 2 Wiki" href={wikiLink} />
                  &nbsp;&middot;&nbsp;
                  <Typography variant="body2" component={Link} href={`${item.slug}/#reserved`}>
                    {item.reserved_count} Reserved
                  </Typography>
                  &nbsp;&middot;&nbsp;
                  <Typography variant="body2" component={Link} href={`${item.slug}/#delivered`}>
                    {item.sold_count} Delivered
                  </Typography>
                  {/* <br /> */}
                  {/* <Typography color="textSecondary" component="span"> */}
                  {/*  {`Median Ask: `} */}
                  {/* </Typography> */}
                  {/* {item.median_ask.toFixed(2)} */}
                </Typography>
                <BidButton
                  onClick={handleBuyOrderClick}
                  className={classes.postItemButton}
                  style={{ marginTop: 1 }}
                  variant="outlined"
                  disableUnderline
                  fullWidth>
                  Place buy order
                </BidButton>
              </Typography>
            </div>
          ) : (
            /* mobile screen */
            <div>
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

              <Typography
                noWrap
                component="h1"
                variant="h6"
                style={
                  item.rarity !== 'regular' ? { color: itemRarityColorMap[item.rarity] } : null
                }>
                {item.name}
              </Typography>
              <Typography>
                <Link href={`/search?hero=${item.hero}`}>{item.hero}</Link>
              </Typography>
              <Typography
                color="textSecondary"
                variant="body2"
                component={Link}
                href={`/search?origin=${item.origin}`}>
                {item.origin}
                {item.rarity !== 'regular' && (
                  <>
                    &nbsp;&middot;
                    <RarityTag
                      color="textSecondary"
                      variant="body2"
                      component={Link}
                      rarity={item.rarity}
                      href={`/search?rarity=${item.rarity}`}
                    />
                  </>
                )}
              </Typography>
              <div style={{ marginTop: 8 }}>
                <ChipLink label="Dota 2 Wiki" href={wikiLink} />
              </div>

              <br />
              {isLoggedIn && (
                <div align="center" style={{ display: 'flex', marginBottom: 2 }}>
                  <Button
                    className={classes.postItemButton}
                    variant="outlined"
                    color="secondary"
                    component={Link}
                    href={`/post-item?s=${item.slug}`}
                    disableUnderline>
                    Post this Item
                  </Button>
                  <BidButton className={classes.postItemButton} variant="outlined" disableUnderline>
                    Place Buy Order
                  </BidButton>
                </div>
              )}
            </div>
          )}

          <MarketList
            offers={markets}
            buyOrders={buyOrders}
            error={error}
            pagination={
              !error && (
                <TablePaginationRouter
                  linkProps={linkProps}
                  style={{ textAlign: 'right' }}
                  count={markets.total_count || 0}
                  page={filter.page}
                />
              )
            }
          />
        </Container>
        <BuyOrderDialog
          catalog={item}
          open={openBuyOrderDialog}
          onClose={() => {
            setOpenBuyOrderDialog(false)
          }}
        />
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
