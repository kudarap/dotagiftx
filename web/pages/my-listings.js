import React from 'react'
import useSWR from 'swr'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { fetcherWithToken, MY_MARKETS } from '@/service/api'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import MyMarketList from '@/components/MyMarketList'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(6),
  },
}))

const activeMarketFilter = {
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
}

export default function About() {
  const classes = useStyles()

  const { data: activeListings, activeListingsError } = useSWR(
    [MY_MARKETS, activeMarketFilter],
    fetcherWithToken
  )

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            My Active Listings
          </Typography>

          {activeListingsError && <div>failed to load active listings</div>}
          {!activeListings && <LinearProgress color="secondary" />}
          {!activeListingsError && activeListings && <MyMarketList datatable={activeListings} />}
        </Container>
      </main>

      <Footer />
    </>
  )
}
