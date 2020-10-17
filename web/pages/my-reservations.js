import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import * as format from '@/lib/format'
import { myMarketSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { MARKET_STATUS_RESERVED } from '@/constants/market'
import ReservationList from '@/components/ReservationList'
import TablePagination from '@/components/TablePagination'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
}))

const marketFilter = {
  status: MARKET_STATUS_RESERVED,
  sort: 'updated_at:desc',
  page: 1,
}

const initialDatatable = {
  data: [],
  result_count: 0,
  total_count: 0,
  loading: false,
  error: null,
}

export default function MyReservations() {
  const classes = useStyles()

  const [data, setData] = React.useState(initialDatatable)
  const [total, setTotal] = React.useState(0)
  const [filter, setFilter] = React.useState(marketFilter)

  React.useEffect(() => {
    ;(async () => {
      setData({ ...data, loading: true })
      try {
        const res = await myMarketSearch(filter)
        setData({ ...data, loading: false, ...res })
      } catch (e) {
        setData({ ...data, loading: false, error: e.message })
      }
    })()
  }, [filter])

  React.useEffect(() => {
    ;(async () => {
      const res = await myMarketSearch(filter)
      setTotal(res.total_count)
    })()
  }, [])

  const handlePageChange = (e, page) => {
    setFilter({ ...filter, page })
  }

  const handleSearchInput = value => {
    setFilter({ ...filter, loading: true, page: 1, q: value })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography component="h1" gutterBottom>
            Buyer Reservations {total !== 0 && `(${format.numberWithCommas(total)})`}
          </Typography>

          <ReservationList
            datatable={data}
            loading={data.loading}
            error={data.error}
            onSearchInput={handleSearchInput}
          />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={data.total_count || 0}
            page={filter.page}
            onChangePage={handlePageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
