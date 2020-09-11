import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { myMarketSearch } from '@/service/api'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import MyMarketList from '@/components/MyMarketList'
import TablePagination from '@/components/TablePagination'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(6),
  },
}))

const activeMarketFilter = {
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
  page: 1,
}

const initialDatatable = {
  data: [],
  result_count: 0,
  total_count: 0,
  loading: true,
  error: null,
}

export default function MyListings() {
  const classes = useStyles()

  const [activeLists, setActiveLists] = React.useState(initialDatatable)
  const [filter, setFilter] = React.useState(activeMarketFilter)

  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await myMarketSearch(filter)
        setActiveLists({ ...activeLists, loading: false, ...res })
      } catch (e) {
        setActiveLists({ ...activeLists, loading: false, error: e.message })
      }
    })()
  }, [filter])

  const handlePageChange = (e, page) => {
    setFilter({ ...filter, page })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            My Active Listings
          </Typography>

          {activeLists.error && <div>failed to load active listings</div>}
          {activeLists.loading && <LinearProgress color="secondary" />}
          <MyMarketList datatable={activeLists} />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={activeLists.total_count || 0}
            page={filter.page}
            onChangePage={handlePageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
