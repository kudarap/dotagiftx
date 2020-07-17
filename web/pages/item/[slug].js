import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { catalog, CDN_URL, marketSearch, trackViewURL } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import RarityTag from '@/components/RarityTag'
import MarketList from '@/components/MarketList'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import TablePagination from '@/components/TablePaginationRouter'

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
}))

export default function ItemDetails({ item, markets, canonicalURL }) {
  const classes = useStyles()

  const router = useRouter()
  const [page, setPage] = React.useState(Number(router.query.page || 1))

  const handlePageChange = (e, p) => {
    setPage(p)
  }

  const linkProps = { href: '/item/[slug]', as: `/item/${item.slug}` }

  const metaTitle = `DotagiftX :: Listings for ${item.name} :: Price starts at $${item.lowest_ask}`
  const metaDesc = `Buy ${item.name} from ${
    item.origin
  } ${item.rarity.toString().toUpperCase()} for ${item.hero}. Price start at ${item.lowest_ask}`

  return (
    <>
      <Head>
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />

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
              <ItemImage
                className={classes.media}
                image={`/300x170/${item.image}`}
                title={item.name}
                rarity={item.rarity}
              />
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

          <MarketList data={markets} />
          <TablePagination
            linkProps={linkProps}
            style={{ textAlign: 'right' }}
            count={markets.total_count}
            page={page}
            onChangePage={handlePageChange}
          />
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
  markets: PropTypes.object,
}
ItemDetails.defaultProps = {
  markets: {
    data: [],
  },
}

const marketSearchFilter = { status: MARKET_STATUS_LIVE, sort: 'price', page: 1 }

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query, req } = props
  const item = await catalog(params.slug)

  const filter = { ...marketSearchFilter, item_id: item.id }
  if (query.page) {
    filter.page = Number(query.page)
  }

  const canonicalURL = req.headers.host + req.url

  return {
    props: {
      item,
      markets: await marketSearch(filter),
      canonicalURL,
    },
  }
}
