import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { myMarketSearch } from '@/service/api'
import { MARKET_STATUS_CANCELLED, MARKET_STATUS_SOLD } from '@/constants/market'
import HistoryList from '@/components/HistoryList'
import TablePagination from '@/components/TablePagination'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(6),
  },
}))

const activeMarketFilter = {
  status: MARKET_STATUS_SOLD,
  sort: 'updated_at:desc',
  page: 1,
}
const cancelledMarketFilter = {
  status: MARKET_STATUS_CANCELLED,
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

export default function History() {
  const classes = useStyles()

  const [soldItems, setSoldItems] = React.useState(initialDatatable)
  const [soldFilter, setSoldFilter] = React.useState(activeMarketFilter)

  const [cancelledItems, setCancelledItems] = React.useState(initialDatatable)
  const [cancelledFilter, setCancelledFilter] = React.useState(cancelledMarketFilter)

  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await myMarketSearch(soldFilter)
        setSoldItems({ ...soldItems, loading: false, ...res })
      } catch (e) {
        setSoldItems({ ...soldItems, loading: false, error: e.message })
      }

      try {
        const res = await myMarketSearch(cancelledFilter)
        setCancelledItems({ ...cancelledItems, loading: false, ...res })
      } catch (e) {
        setCancelledItems({ ...cancelledItems, loading: false, error: e.message })
      }
    })()
  }, [soldFilter])

  const handleSoldPageChange = (e, page) => {
    setSoldFilter({ ...soldFilter, page })
  }
  const handleCancelledPageChange = (e, page) => {
    setCancelledFilter({ ...soldFilter, page })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography id="delivered" component="h1" gutterBottom>
            Delivered Items
          </Typography>
          {soldItems.error && <div>failed to load sold items</div>}
          {soldItems.loading && <LinearProgress color="secondary" />}
          <HistoryList datatable={soldItems} />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={soldItems.total_count || 0}
            page={soldFilter.page}
            onChangePage={handleSoldPageChange}
          />
          <br />

          <Typography id="cancelled" component="h1" gutterBottom>
            Cancelled Items
          </Typography>
          {cancelledItems.error && <div>failed to load cancelled</div>}
          {cancelledItems.loading && <LinearProgress color="secondary" />}
          <HistoryList datatable={cancelledItems} />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={cancelledItems.total_count || 0}
            page={cancelledFilter.page}
            onChangePage={handleCancelledPageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
