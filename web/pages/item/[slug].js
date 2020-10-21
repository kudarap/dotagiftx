import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import Typography from '@material-ui/core/Typography'
import { get as getLoggedInUser, isOk as checkLoggedIn } from '@/service/auth'
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
import { itemRarityColorMap } from '@/constants/palette'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(0),
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
      margin: '0 auto !important',
      width: 300,
      height: 170,
    },
    width: 164,
    height: 109,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(1.5),
  },
  postItemButton: {
    [theme.breakpoints.down('xs')]: {
      margin: `8px auto !important`,
      width: 300,
    },
    width: 164,
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
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))

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
  const [error, setError] = React.useState(null)

  // Handle market request on page change.
  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await marketSearch(filter)
        setMarkets(res)
      } catch (e) {
        setError(e.message)
      }
    })()
  }, [filter])

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

  const isLoggedIn = checkLoggedIn()
  let currentUserID = null
  if (isLoggedIn) {
    currentUserID = getLoggedInUser().user_id
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
        <meta name="twitter:site" content="@DotagiftX" />
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
                      image={`/300x170/${item.image}`}
                      title={item.name}
                      rarity={item.rarity}
                    />
                  </a>
                  {isLoggedIn && (
                    <Button
                      className={classes.postItemButton}
                      variant="outlined"
                      color="secondary"
                      size="small"
                      component={Link}
                      href={`/post-item?s=${item.slug}`}
                      disableUnderline
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
                      <RarityTag
                        rarity={item.rarity}
                        variant="body1"
                        component={Link}
                        href={`/search?q=${item.rarity}`}
                      />
                    </>
                  )}
                  <br />
                  <Typography color="textSecondary" component="span">
                    {`Used by: `}
                  </Typography>
                  <Link href={`/search?q=${item.hero}`}>{item.hero}</Link>
                  <br />
                  <Typography color="textSecondary" component="span">
                    {`Links: `}
                  </Typography>
                  <ChipLink label="Dota 2 Wiki" href={wikiLink} />
                </Typography>
              </Typography>
            </div>
          ) : (
            /* mobile screen */
            <div>
              {item.image && (
                <a href={wikiLink} target="_blank" rel="noreferrer noopener">
                  <ItemImage
                    className={classes.media}
                    image={`/300x170/${item.image}`}
                    title={item.name}
                  />
                </a>
              )}
              <div align="center">
                <Button
                  className={classes.postItemButton}
                  variant="outlined"
                  color="secondary"
                  size="small"
                  component={Link}
                  href={`/post-item?s=${item.slug}`}
                  disableUnderline
                  fullWidth>
                  Post this Item
                </Button>
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
                <Link href={`/search?q=${item.hero}`}>{item.hero}</Link>
              </Typography>
              <Typography
                color="textSecondary"
                variant="body2"
                component={Link}
                href={`/search?q=${item.origin}`}>
                {item.origin}
                {item.rarity !== 'regular' && (
                  <>
                    &nbsp;&middot;
                    <RarityTag
                      color="textSecondary"
                      variant="body2"
                      component={Link}
                      rarity={item.rarity}
                      href={`/search?q=${item.rarity}`}
                    />
                  </>
                )}
              </Typography>
              <div style={{ marginTop: 8 }}>
                <ChipLink label="Dota 2 Wiki" href={wikiLink} />
              </div>
            </div>
          )}
          <br />

          <MarketList data={markets} currentUserID={currentUserID} error={error} />
          {!error && (
            <TablePaginationRouter
              linkProps={linkProps}
              style={{ textAlign: 'right' }}
              count={markets.total_count}
              page={filter.page}
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

// This gets called on every request
export async function getServerSideProps(props) {
  const { res, params } = props
  res.setHeader('location', `/${params.slug}`)
  res.statusCode = 302
  res.end()
  return { props: {} }
}
