import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { myMarketSearch } from '@/service/api'
import { MARKET_STATUS_RESERVED, MARKET_STATUS_SOLD } from '@/constants/market'
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

  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await myMarketSearch(soldFilter)
        setSoldItems({ ...soldItems, loading: false, ...res })
      } catch (e) {
        setSoldItems({ ...soldItems, loading: false, error: e.message })
      }
    })()
  }, [soldFilter])

  const handlePageChange = (e, page) => {
    setSoldFilter({ ...soldFilter, page })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Delivered Items
          </Typography>

          {soldItems.error && <div>failed to load soldItems</div>}
          {soldItems.loading && <LinearProgress color="secondary" />}
          <HistoryList datatable={soldItems} />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={soldItems.total_count || 0}
            page={soldFilter.page}
            onChangePage={handlePageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
