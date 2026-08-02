import React from 'react'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import LinearProgress from '@mui/material/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { myMarketSearch } from '@/service/api'
import {
  MARKET_STATUS_BID_COMPLETED,
  MARKET_STATUS_CANCELLED,
  MARKET_STATUS_SOLD,
} from '@/constants/market'
import HistoryList from '@/components/HistoryList'
import TablePagination from '@/components/TablePagination'
import type { MarketDatatable } from '@/components/MarketList'
import type { QueryFilter } from '@/service/http'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
}))

const activeMarketFilter = {
  status: MARKET_STATUS_SOLD,
  sort: 'updated_at:desc',
  page: 1,
}
const completedBidFilter = {
  status: MARKET_STATUS_BID_COMPLETED,
  sort: 'updated_at:desc',
  page: 1,
}
const cancelledMarketFilter = {
  status: MARKET_STATUS_CANCELLED,
  sort: 'updated_at:desc',
  page: 1,
}

interface Datatable extends MarketDatatable {
  result_count?: number
  loading: boolean
  error: string | null
}

const initialDatatable: Datatable = {
  data: [],
  result_count: 0,
  total_count: 0,
  loading: true,
  error: null,
}

export default function MyHistory() {
  const { classes } = useStyles()

  const [soldItems, setSoldItems] = React.useState<Datatable>(initialDatatable)
  const [soldFilter, setSoldFilter] = React.useState<QueryFilter & { page?: number }>(
    activeMarketFilter
  )

  const [completedItems, setCompletedItems] = React.useState<Datatable>(initialDatatable)
  const [completedFilter, setCompletedFilter] = React.useState<QueryFilter & { page?: number }>(
    completedBidFilter
  )

  const [cancelledItems, setCancelledItems] = React.useState<Datatable>(initialDatatable)
  const [cancelledFilter, setCancelledFilter] = React.useState<QueryFilter & { page?: number }>(
    cancelledMarketFilter
  )

  React.useEffect(() => {
    let active = true
    ;(async () => {
      try {
        const res = (await myMarketSearch(soldFilter)) as MarketDatatable
        if (active) {
          setSoldItems(current => ({ ...current, loading: false, ...res }))
        }
      } catch (e) {
        if (active) {
          setSoldItems(current => ({ ...current, loading: false, error: (e as Error).message }))
        }
      }
    })()

    return () => {
      active = false
    }
  }, [soldFilter])

  React.useEffect(() => {
    let active = true
    ;(async () => {
      try {
        const res = (await myMarketSearch(completedFilter)) as MarketDatatable
        if (active) {
          setCompletedItems(current => ({ ...current, loading: false, ...res }))
        }
      } catch (e) {
        if (active) {
          setCompletedItems(current => ({ ...current, loading: false, error: (e as Error).message }))
        }
      }
    })()

    return () => {
      active = false
    }
  }, [completedFilter])

  React.useEffect(() => {
    let active = true
    ;(async () => {
      try {
        const res = (await myMarketSearch(cancelledFilter)) as MarketDatatable
        if (active) {
          setCancelledItems(current => ({ ...current, loading: false, ...res }))
        }
      } catch (e) {
        if (active) {
          setCancelledItems(current => ({ ...current, loading: false, error: (e as Error).message }))
        }
      }
    })()

    return () => {
      active = false
    }
  }, [cancelledFilter])

  const handleSoldPageChange = (e: unknown, page: number) => {
    setSoldFilter({ ...soldFilter, page })
  }
  const handleCompletedPageChange = (e: unknown, page: number) => {
    setCompletedFilter({ ...completedFilter, page })
  }
  const handleCancelledPageChange = (e: unknown, page: number) => {
    setCancelledFilter({ ...cancelledFilter, page })
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
            style={{ textAlign: 'right', minHeight: 48 }}
            count={soldItems.total_count || 0}
            page={soldFilter.page as number}
            onPageChange={handleSoldPageChange}
          />

          <Typography id="delivered" component="h1" gutterBottom>
            Completed Orders
          </Typography>
          {completedItems.error && <div>failed to load completed orders</div>}
          {completedItems.loading && <LinearProgress color="secondary" />}
          <HistoryList datatable={completedItems} />
          <TablePagination
            style={{ textAlign: 'right', minHeight: 48 }}
            count={completedItems.total_count || 0}
            page={completedFilter.page as number}
            onPageChange={handleCompletedPageChange}
          />

          <Typography id="cancelled" component="h1" gutterBottom>
            Cancelled Items
          </Typography>
          {cancelledItems.error && <div>failed to load cancelled</div>}
          {cancelledItems.loading && <LinearProgress color="secondary" />}
          <HistoryList datatable={cancelledItems} />
          <TablePagination
            style={{ textAlign: 'right', minHeight: 48 }}
            count={cancelledItems.total_count || 0}
            page={cancelledFilter.page as number}
            onPageChange={handleCancelledPageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
