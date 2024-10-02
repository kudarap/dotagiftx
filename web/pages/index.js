import React from 'react'
import PropTypes from 'prop-types'
import useSWR from 'swr'
import Head from 'next/head'
import Router from 'next/router'
import dynamic from 'next/dynamic'
import { makeStyles } from 'tss-react/mui'
import LinearProgress from '@mui/material/LinearProgress'
import Typography from '@mui/material/Typography'
import Grid from '@mui/material/Grid'
import Divider from '@mui/material/Divider'
import { APP_NAME, APP_URL } from '@/constants/strings'
import {
  CATALOGS,
  catalogTrendSearch,
  fetcher,
  STATS_TOP_HEROES,
  STATS_TOP_ORIGINS,
  statsMarketSummary,
} from '@/service/api'
import * as format from '@/lib/format'
import Header from '@/components/Header'
import Container from '@/components/Container'
// import SearchInput from '@/components/SearchInput'
// import CatalogList from '@/components/CatalogList'
// import Link from '@/components/Link'
// import Footer from '@/components/Footer'

// const Header = dynamic(() => import('@/components/Header'))
// const Container = dynamic(() => import('@/components/Container'))
const SearchInput = dynamic(() => import('@/components/SearchInput'))
const CatalogList = dynamic(() => import('@/components/CatalogList'))
const Link = dynamic(() => import('@/components/Link'))
const Footer = dynamic(() => import('@/components/Footer'))

const useStyles = makeStyles()(theme => ({
  main: {},
  searchBar: {
    margin: '0 auto',
    marginBottom: theme.spacing(4),
  },
  banner: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(0),
    },
    margin: theme.spacing(0, 0, 2, 0),
    padding: theme.spacing(1.5),
    border: '1px solid #52564e',
    background: '#2d3431',
    borderRadius: 4,
  },
  bannerHighlight: {
    background: '-webkit-linear-gradient(#EBCF87 10%, #EA6953 90%)',
    backgroundClip: 'border-box',
    backgroundClip: 'text',
    WebkitBackgroundClip: 'text',
    WebkitTextFillColor: 'transparent',
    filter: 'drop-shadow(0px 0px 5px #e1261c)',
  },
  bannerText: {
    [theme.breakpoints.down('sm')]: {
      fontSize: theme.typography.body2.fontSize,
    },
  },
  footLinks: {
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
    },
  },
  divider: {
    marginTop: theme.spacing(3),
    marginBottom: theme.spacing(3),
  },
}))

const recentItemsFilter = {
  sort: 'recent',
  limit: 5,
}
const recentBidItemsFilter = {
  sort: 'recent-bid',
  limit: 5,
}
const topSellerItemsFilter = {
  sort: 'sale_count:desc',
}

export default function Index({ marketSummary, trendingItems }) {
  const { classes } = useStyles()

  const { data: recentBidItems, error: recentBidError } = useSWR(
    [CATALOGS, recentBidItemsFilter],
    fetcher
  )
  const { data: recentItems, error: recentError } = useSWR(
    recentBidItems ? [CATALOGS, recentItemsFilter] : null,
    fetcher
  )
  const { data: topSellers } = useSWR(
    recentItems ? [CATALOGS, topSellerItemsFilter] : null,
    fetcher
  )
  const { data: topOrigins } = useSWR(topSellers ? STATS_TOP_ORIGINS : null, fetcher)
  const { data: topHeroes } = useSWR(topOrigins ? STATS_TOP_HEROES : null, fetcher)

  const handleSubmit = keyword => {
    Router.push(`/search?q=${keyword}`)
  }

  const description = `Search on ${marketSummary.live} Giftable items`

  const metaTitle = `${APP_NAME} :: Dota 2 Giftables Community Market`
  const metaDesc = `${description}. ${APP_NAME} was made to provide better search and pricing for 
          Dota 2 Giftable items like Collector's Caches which are not available on Steam Community Market. 
          The project was heavily inspired by Giftable Megathread from r/Dota2Trade.`

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={APP_URL} />

        {/* Twitter Card */}
        <meta name="twitter:card" content="summary" />
        <meta name="twitter:title" content={metaTitle} />
        <meta name="twitter:description" content={metaDesc} />
        <meta name="twitter:image" content={`${APP_URL}/icon.png`} />
        <meta name="twitter:site" content={`@${APP_NAME}`} />
        {/* OpenGraph */}
        <meta property="og:url" content={APP_URL} />
        <meta property="og:type" content="website" />
        <meta property="og:title" content={metaTitle} />
        <meta property="og:description" content={metaDesc} />
        <meta property="og:image" content={`${APP_URL}/icon.png`} />
      </Head>

      <Header disableSearch />

      <main className={classes.main}>
        <div
          style={{
            width: '100%',
            height: 340,
            marginBottom: 500 - 340,
            maskImage: 'linear-gradient(to top, transparent 0%, black 90%)',
            WebkitMaskImage: 'linear-gradient(to top, transparent 0%, black 90%)',
            position: 'relative',
            zIndex: 0,
          }}>
          <div
            style={{
              background: 'url(/assets/ti_ringmaster_banner.png) no-repeat center center',
              // background: 'url(https://cdn.akamai.steamstatic.com/apps/dota2/images/dota_react/international2024/esports_site/footer_bg01.png) no-repeat center center',
              backgroundColor: '#263238',
              backgroundSize: 'cover',
              backgroundPositionY: -120,
              width: '100%',
              height: '100%',
            }}></div>
        </div>

        <Container
          sx={{
            mt: {
              md: -35,
              xs: -61,
            },
            position: 'relative',
          }}>
          <div className={classes.banner} hidden>
            <Typography
              className={classes.bannerText}
              component="h1"
              variant="body2"
              color="textSecondary">
              <Typography className={classes.bannerHighlight} color="secondary" component="span">
                {APP_NAME}
              </Typography>{' '}
              was made to provide better search and pricing for Dota 2 Giftable items like
              Collector&apos;s Caches which are not available on{' '}
              <Link href="https://steamcommunity.com" rel="noreferrer noopener" target="_blank">
                Steam Community Market
              </Link>
              . The project was heavily inspired by <strong>Giftable Megathread</strong> from{' '}
              <Link
                href="https://www.reddit.com/r/Dota2Trade"
                rel="noreferrer noopener"
                target="_blank">
                r/Dota2Trade
              </Link>
              .
            </Typography>
            {/* <Typography className={classes.bannerText} variant="h3" component="h1" align="center"> */}
            {/*  /!* Search for Dota 2 <span style={{ display: 'inline-block' }}>Giftable items</span> *!/ */}
            {/*  /!* Buy & Sell *!/ */}
            {/*  Search for <span style={{ display: 'inline-block' }}>Dota 2 giftabe items</span> */}
            {/* </Typography> */}
          </div>

          <SearchInput label={description} onSubmit={handleSubmit} />
          {/* <Typography variant="caption" sx={{ mt: -2.5, mr: 2.5, float: 'right' }}>
            <Link href="/giveaway" target="_blank" rel="noreferrer noopener" color="secondary">
              Collector's Cache Giveaway
            </Link>
          </Typography> */}
          <br />

          {/* Trending Items */}
          <Typography>Trending</Typography>
          {trendingItems.error && <div>failed to load trending items: {trendingItems.error}</div>}
          {!trendingItems.error && <CatalogList items={trendingItems.data} />}
          <br />

          {/* Recent Buy Orders */}
          <Typography>
            New Buy Orders
            <Link
              href={`/search?sort=${recentBidItemsFilter.sort}`}
              color="secondary"
              style={{ float: 'right' }}>
              See All
            </Link>
          </Typography>
          {recentBidError && <div>failed to load recent buy orders items</div>}
          {!recentBidItems && <LinearProgress color="secondary" />}
          {!recentBidError && recentBidItems && (
            <CatalogList items={recentBidItems.data} variant="recent" bidType />
          )}
          <br />

          {/* Recent Market items */}
          <Typography>
            New Sell Listings
            <Link
              href={`/search?sort=${recentItemsFilter.sort}`}
              color="secondary"
              style={{ float: 'right' }}>
              See All
            </Link>
          </Typography>
          {recentError && <div>failed to load recent items</div>}
          {!recentItems && <LinearProgress color="secondary" />}
          {!recentError && recentItems && <CatalogList items={recentItems.data} variant="recent" />}
          <br />

          {/* Market stats */}
          <Divider className={classes.divider} light variant="middle" />
          <Grid container spacing={2} style={{ textAlign: 'center' }}>
            <Grid item sm={3} xs={6} component={Link} href="/search" disableUnderline>
              <Typography variant="h4" component="span">
                {marketSummary.live}
              </Typography>
              <br />
              <Typography color="textSecondary" variant="body2">
                <em>Available Offers</em>
              </Typography>
            </Grid>
            <Grid
              item
              sm={3}
              xs={6}
              component={Link}
              href="/search?sort=recent-bid"
              disableUnderline>
              <Typography variant="h4" component="span">
                {marketSummary.bids.bid_live}
              </Typography>
              <br />
              <Typography color="textSecondary" variant="body2">
                <em>Buy Orders</em>
              </Typography>
            </Grid>
            <Grid item sm={3} xs={6} component={Link} href="/history/reserved" disableUnderline>
              <Typography variant="h4" component="span">
                {marketSummary.reserved}
              </Typography>
              <br />
              <Typography color="textSecondary" variant="body2">
                <em>On Reserved</em>
              </Typography>
            </Grid>
            <Grid item sm={3} xs={6} component={Link} href="/history/delivered" disableUnderline>
              <Typography variant="h4" component="span">
                {marketSummary.sold}
              </Typography>
              <br />
              <Typography color="textSecondary" variant="body2">
                <em>Delivered Items</em>
              </Typography>
            </Grid>
          </Grid>
          <Divider className={classes.divider} light variant="middle" />
          <br />

          {/* Top 10 foot links */}
          <Grid container spacing={2}>
            {/* Top 10 Heroes*/}
            <Grid item sm={4} xs={12}>
              <Typography className={classes.footLinks}>Top Heroes</Typography>
              {topHeroes &&
                topHeroes.map(hero => (
                  <Link
                    key={hero}
                    href={`/search?hero=${hero}`}
                    color="secondary"
                    className={classes.footLinks}>
                    <Typography variant="subtitle1" component="p">
                      {hero}
                    </Typography>
                  </Link>
                ))}
            </Grid>
            {/* Top 10 Sellers */}
            <Grid item sm={4} xs={12}>
              <Typography className={classes.footLinks}>Top Sellers</Typography>
              {topSellers &&
                topSellers.data.map(item => (
                  <Link
                    key={item.slug}
                    href={`/${item.slug}`}
                    color="secondary"
                    className={classes.footLinks}>
                    <Typography variant="subtitle1" component="p">
                      {item.name}
                    </Typography>
                  </Link>
                ))}
            </Grid>
            {/* Top 10 Treasures */}
            <Grid item sm={4} xs={12}>
              <Typography className={classes.footLinks}>Top Treasures</Typography>
              {topOrigins &&
                topOrigins.map(origin => (
                  <Link
                    key={origin}
                    href={`/search?origin=${origin}`}
                    color="secondary"
                    className={classes.footLinks}>
                    <Typography variant="subtitle1" component="p">
                      {origin}
                    </Typography>
                  </Link>
                ))}
            </Grid>
          </Grid>
        </Container>
      </main>

      <Footer />
    </>
  )
}
Index.propTypes = {
  marketSummary: PropTypes.object.isRequired,
  trendingItems: PropTypes.object.isRequired,
}

// This gets called on every request
export async function getServerSideProps() {
  const marketSummary = await statsMarketSummary()
  marketSummary.live = format.numberWithCommas(marketSummary.live)
  marketSummary.reserved = format.numberWithCommas(marketSummary.reserved)
  marketSummary.sold = format.numberWithCommas(marketSummary.sold)
  marketSummary.bids.bid_live = format.numberWithCommas(marketSummary.bids.bid_live)

  let trendingItems = { error: null }
  try {
    trendingItems = await catalogTrendSearch()
  } catch (e) {
    trendingItems.error = e
  }

  return {
    props: {
      marketSummary,
      trendingItems,
      unstable_revalidate: 60,
    },
  }
}
