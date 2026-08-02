import React from 'react'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import * as format from '@/lib/format'
import { myMarketSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { MARKET_STATUS_RESERVED } from '@/constants/market'
import ReservationList from '@/components/ReservationList'
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

const marketFilter = {
  status: MARKET_STATUS_RESERVED,
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
  loading: false,
  error: null,
}

export default function MyReservations() {
  const { classes } = useStyles()

  const [data, setData] = React.useState<Datatable>(initialDatatable)
  const [total, setTotal] = React.useState(0)
  const [filter, setFilter] = React.useState<QueryFilter & { page?: number; q?: string }>(
    marketFilter
  )
  const [reloadFlag, setReloadFlag] = React.useState(false)

  React.useEffect(() => {
    ;(async () => {
      setData(current => ({ ...current, loading: true, error: null }))
      try {
        const res = (await myMarketSearch(filter)) as MarketDatatable
        setData(current => ({ ...current, loading: false, ...res }))
      } catch (e) {
        setData(current => ({ ...current, loading: false, error: (e as Error).message }))
      }
    })()
  }, [filter, reloadFlag])

  React.useEffect(() => {
    ;(async () => {
      const res = (await myMarketSearch(filter)) as MarketDatatable
      setTotal(res.total_count)
    })()
  }, [filter])

  const handleSearchInput = (value: string) => {
    setFilter({ ...filter, page: 1, q: value })
  }
  const handlePageChange = (e: unknown, page: number) => {
    setFilter({ ...filter, page })
  }
  const handleReloadToggle = () => {
    setReloadFlag(!reloadFlag)
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
            onReload={handleReloadToggle}
          />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={data.total_count || 0}
            page={filter.page as number}
            onPageChange={handlePageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
