import React from 'react'
import PropTypes from 'prop-types'
import useSWR from 'swr'
import Head from 'next/head'
import Router from 'next/router'
import dynamic from 'next/dynamic'
import { makeStyles } from '@material-ui/core/styles'
import LinearProgress from '@material-ui/core/LinearProgress'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Divider from '@material-ui/core/Divider'
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

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  searchBar: {
    margin: '0 auto',
    marginBottom: theme.spacing(4),
  },
  banner: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(0),
    },
    margin: theme.spacing(0, 0, 2, 0),
    padding: theme.spacing(1.5),
    border: '1px solid #52564e',
    background: '#2d3431',
    borderRadius: 4,
  },
  bannerText: {
    [theme.breakpoints.down('xs')]: {
      fontSize: theme.typography.body2.fontSize,
    },
  },
  footLinks: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
    },
  },
  divider: {
    marginTop: theme.spacing(3),
    marginBottom: theme.spacing(3),
  },
}))

const popularItemsFilter = {
  sort: 'popular',
  limit: 5,
}
const recentItemsFilter = {
  sort: 'recent',
  limit: 5,
}
const topSellerItemsFilter = {
  sort: 'sale_count:desc',
}

export default function Index({ marketSummary, trendingItems }) {
  const classes = useStyles()

  const { data: recentItems, recentError } = useSWR([CATALOGS, recentItemsFilter], fetcher)
  const { data: popularItems, popularError } = useSWR([CATALOGS, popularItemsFilter], fetcher)
  const { data: topSellers } = useSWR([CATALOGS, topSellerItemsFilter], fetcher)
  const { data: topOrigins } = useSWR(STATS_TOP_ORIGINS, fetcher)
  const { data: topHeroes } = useSWR(STATS_TOP_HEROES, fetcher)

  const handleSubmit = keyword => {
    Router.push(`/search?q=${keyword}`)
  }

  const description = `Search on ${marketSummary.live} giftable items`

  const metaTitle = `${APP_NAME} :: Dota 2 Giftables Community Market`
  const metaDesc = `${description}. ${APP_NAME} was made to provide better search and pricing for 
          Dota 2 Giftable items like Collector's Caches which are not available on Steam Community Market. 
          The project was heavily inspired by Giftable Megathread from r/Dota2Trade.`

  return (
    <>
      <Head>
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
        <Container>
          <div className={classes.banner}>
            <Typography
              className={classes.bannerText}
              component="h1"
              variant="body2"
              color="textSecondary">
              <Typography color="secondary" component="span">
                {APP_NAME}
              </Typography>{' '}
              was made to provide better search and pricing for Dota 2 giftable items like
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
          <br />

          {/* Trending Items */}
          <Typography>Trending Items</Typography>
          {trendingItems.error && <div>failed to load popular items: {trendingItems.error}</div>}
          {!trendingItems.error && <CatalogList items={trendingItems.data} />}
          <br />

          {/* Recent Market items */}
          <Typography>
            Recently Posted
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

          {/* Popular Market items */}
          <Typography>
            Most Popular
            <Link
              href={`/search?sort=${popularItemsFilter.sort}`}
              color="secondary"
              style={{ float: 'right' }}>
              See All
            </Link>
          </Typography>
          {popularError && <div>failed to load popular items</div>}
          {!popularItems && <LinearProgress color="secondary" />}
          {!popularError && popularItems && <CatalogList items={popularItems.data} />}
          <br />

          {/* Market stats */}
          <Divider className={classes.divider} light variant="middle" />
          <Grid container spacing={2} style={{ textAlign: 'center' }}>
            <Grid item sm={4} xs={12} component={Link} href="/search" disableUnderline>
              <Typography variant="h4" component="span">
                {marketSummary.live}
              </Typography>
              <br />
              <Typography color="textSecondary" variant="body2">
                <em>Available Offers</em>
              </Typography>
            </Grid>
            <Grid item sm={4} xs={6} component={Link} href="/history?reserved" disableUnderline>
              <Typography variant="h4" component="span">
                {marketSummary.reserved}
              </Typography>
              <br />
              <Typography color="textSecondary" variant="body2">
                <em>On Reserved</em>
              </Typography>
            </Grid>
            <Grid item sm={4} xs={6} component={Link} href="/history?delivered" disableUnderline>
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
            <Grid item sm={4} xs={12}>
              <Typography className={classes.footLinks}>Top Sellers</Typography>
              {topSellers &&
                topSellers.data.map(item => (
                  <Link href={`/${item.slug}`} color="secondary" className={classes.footLinks}>
                    <Typography variant="subtitle1" component="p">
                      {item.name}
                    </Typography>
                  </Link>
                ))}
            </Grid>
            <Grid item sm={4} xs={12}>
              <Typography className={classes.footLinks}>Top Heroes</Typography>
              {topHeroes &&
                topHeroes.map(hero => (
                  <Link
                    href={`/search?hero=${hero}`}
                    color="secondary"
                    className={classes.footLinks}>
                    <Typography variant="subtitle1" component="p">
                      {hero}
                    </Typography>
                  </Link>
                ))}
            </Grid>
            <Grid item sm={4} xs={12}>
              <Typography className={classes.footLinks}>Top Treasures</Typography>
              {topOrigins &&
                topOrigins.map(origin => (
                  <Link
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
