import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { myMarketSearch } from '@/service/api'
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
  loading: true,
  error: null,
}

export default function Reservations() {
  const classes = useStyles()

  const [reservations, setReservations] = React.useState(initialDatatable)
  const [filter, setFilter] = React.useState(marketFilter)

  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await myMarketSearch(filter)
        setReservations({ ...reservations, loading: false, ...res })
      } catch (e) {
        setReservations({ ...reservations, loading: false, error: e.message })
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
          <Typography component="h1" gutterBottom>
            Buyer Reservations
          </Typography>

          {reservations.error && <div>failed to load reservations</div>}
          {reservations.loading && <LinearProgress color="secondary" />}
          <ReservationList datatable={reservations} />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={reservations.total_count || 0}
            page={filter.page}
            onChangePage={handlePageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
