import React from 'react'
import useSWR from 'swr'
import Router from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import LinearProgress from '@material-ui/core/LinearProgress'
import Typography from '@material-ui/core/Typography'
import { CATALOGS, MARKETS, fetcher } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import SearchInput from '@/components/SearchInput'
import CatalogList from '@/components/CatalogList'
import Link from '@/components/Link'

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
    margin: theme.spacing(20, 0, 4, 0),
  },
  bannerText: {
    [theme.breakpoints.down('sm')]: {
      fontSize: 35,
    },
    fontWeight: 'bold',
    color: theme.palette.app.white,
  },
}))

function numberWithCommas(x) {
  return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

function Banner() {
  const classes = useStyles()

  return (
    <div className={classes.banner}>
      <Typography className={classes.bannerText} variant="h3" align="center">
        Search for Dota 2 <span style={{ display: 'inline-block' }}>giftable items</span>
      </Typography>
    </div>
  )
}

const popularItemsFilter = {
  sort: 'view_count:desc',
  limit: 5,
}
const recentItemsFilter = {
  sort: 'recent_ask:desc',
  limit: 5,
}
const marketItemsFilter = {
  limit: 1,
}

export default function Index() {
  const classes = useStyles()

  const { data: popularItems, popularError } = useSWR([CATALOGS, popularItemsFilter], fetcher)
  const { data: recentItems, recentError } = useSWR([CATALOGS, recentItemsFilter], fetcher)
  const { data: marketItems } = useSWR([MARKETS, marketItemsFilter], fetcher)

  const handleSubmit = keyword => {
    Router.push(`/search?q=${keyword}`)
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Banner />

          <SearchInput
            helperText={`Search on ${
              marketItems && numberWithCommas(marketItems.total_count)
            } for sale items`}
            onSubmit={handleSubmit}
          />
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
          {popularError && <div>failed to load</div>}
          {!popularItems && <LinearProgress color="secondary" />}
          {!popularError && popularItems && <CatalogList items={popularItems.data} />}
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
          {recentError && <div>failed to load</div>}
          {!recentItems && <LinearProgress color="secondary" />}
          {!recentError && recentItems && <CatalogList items={recentItems.data} variant="recent" />}
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
