import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import * as format from '@/lib/format'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { myMarketSearch } from '@/service/api'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import MyMarketList from '@/components/MyMarketList'
import TablePagination from '@/components/TablePagination'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  searchPaper: {
    [theme.breakpoints.down('xs')]: {
      width: '100%',
    },
    float: 'right',
    flexGrow: 1,
    padding: '3px 12px',
    marginBottom: 2,
    display: 'flex',
    alignItems: 'center',
    fontStyle: 'italic',
    backgroundColor: theme.palette.primary.dark,
    opacity: 0.8,
  },
  searchInput: {
    marginLeft: theme.spacing(1),
  },
  searchIcon: {
    color: theme.palette.grey[500],
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
  loading: false,
  error: null,
}

export default function MyListings() {
  const classes = useStyles()

  const [data, setData] = React.useState(initialDatatable)
  const [total, setTotal] = React.useState(0)
  const [filter, setFilter] = React.useState(activeMarketFilter)

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
    setFilter({ ...filter, loading: true, q: value })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography component="h1" gutterBottom>
            Active Listings {total !== 0 && `(${format.numberWithCommas(total)})`}
          </Typography>

          <MyMarketList
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
