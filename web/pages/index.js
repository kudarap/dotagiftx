import React from 'react'
import PropTypes from 'prop-types'
import useSWR from 'swr'
import Head from 'next/head'
import Router from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import LinearProgress from '@material-ui/core/LinearProgress'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import {
  CATALOGS,
  fetcher,
  marketSearch,
  catalogSearch,
  STATS_TOP_ORIGINS,
  STATS_TOP_HEROES,
} from '@/service/api'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import * as format from '@/lib/format'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import SearchInput from '@/components/SearchInput'
import CatalogList from '@/components/CatalogList'
import Link from '@/components/Link'
import { APP_URL } from '@/constants/strings'
import Button from '@/components/Button'
import ChipLink from '@/components/ChipLink'

const useStyles = makeStyles(theme => ({
  main: {
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
    margin: theme.spacing(4, 0, 4, 0),
  },
  bannerText: {
    [theme.breakpoints.down('sm')]: {
      fontSize: 35,
    },
    fontWeight: 'bold',
    color: theme.palette.app.white,
  },
  footLinks: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
    },
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

export default function Index({ totalEntries, popularItems }) {
  const classes = useStyles()

  const { data: recentItems, recentError } = useSWR([CATALOGS, recentItemsFilter], fetcher)
  const { data: topOrigins } = useSWR(STATS_TOP_ORIGINS, fetcher)
  const { data: topHeroes } = useSWR(STATS_TOP_HEROES, fetcher)

  const handleSubmit = keyword => {
    Router.push(`/search?q=${keyword}`)
  }

  const description = `Search on ${totalEntries || ''} giftable items`

  const metaTitle = 'DotagiftX - Dota 2 giftable items marketplace'
  const metaDesc = `${description}. DotagiftX was made to provide better search and pricing for 
          Dota 2 giftable items like Collector's Caches which are not available on Steam Community Market. 
          The project was heavily inspired by All Giftable Megathread from r/Dota2Trade.`

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
        <meta name="twitter:image" content={`${APP_URL}/assets/gift.png`} />
        <meta name="twitter:site" content="@DotagiftX" />
        {/* OpenGraph */}
        <meta property="og:url" content={APP_URL} />
        <meta property="og:type" content="website" />
        <meta property="og:title" content={metaTitle} />
        <meta property="og:description" content={metaDesc} />
        <meta property="og:image" content={`${APP_URL}/assets/gift.png`} />
      </Head>

      <Header disableSearch />

      <main className={classes.main}>
        <Container>
          <div className={classes.banner}>
            <Typography component="h1" color="textSecondary">
              <Typography color="secondary" component="span">
                DotagiftX
              </Typography>{' '}
              was made to provide better search and pricing for Dota 2 giftable items like
              Collector&apos;s Caches which are not available on{' '}
              <Link href="https://steamcommunity.com" rel="noreferrer noopener" target="_blank">
                Steam Community Market
              </Link>
              . The project was heavily inspired by <strong>All Giftable Megathread</strong> from{' '}
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

          <SearchInput helperText={description} onSubmit={handleSubmit} />
          <br />

          <Typography>
            Popular Items
            <Link
              href={`/search?sort=${popularItemsFilter.sort}`}
              color="secondary"
              style={{ float: 'right' }}>
              See All
            </Link>
          </Typography>
          {popularItems.error && <div>failed to load popular items: {popularItems.error}</div>}
          {!popularItems.error && <CatalogList items={popularItems.data} />}
          <br />

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

          <Grid container spacing={2}>
            <Grid item sm={6} xs={12}>
              <Typography className={classes.footLinks}>Top Heroes</Typography>
              {topHeroes &&
                topHeroes.map(hero => (
                  <Link href={`/search?q=${hero}`} color="secondary" className={classes.footLinks}>
                    <Typography variant="subtitle1">{hero}</Typography>
                  </Link>
                ))}
            </Grid>
            <Grid item sm={6} xs={12}>
              <Typography className={classes.footLinks}>Top Treasures</Typography>
              {topOrigins &&
                topOrigins.map(origin => (
                  <Link
                    href={`/search?q=${origin}`}
                    color="secondary"
                    className={classes.footLinks}>
                    <Typography variant="subtitle1">{origin}</Typography>
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
  totalEntries: PropTypes.string.isRequired,
  popularItems: PropTypes.object.isRequired,
}

// This gets called on every request
export async function getServerSideProps() {
  const res = await marketSearch({ limit: 1, status: MARKET_STATUS_LIVE })
  const totalEntries = format.numberWithCommas(res.total_count || 0)

  let popularItems = { error: null }
  try {
    popularItems = await catalogSearch(popularItemsFilter)
  } catch (e) {
    popularItems.error = e
  }

  return {
    props: {
      totalEntries,
      popularItems,
      unstable_revalidate: 60,
    },
  }
}
