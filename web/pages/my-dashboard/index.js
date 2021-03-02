import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import { myMarketSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { MARKET_STATUS_LIVE, MARKET_TYPE_ASK } from '@/constants/market'
import MyMarketList from '@/components/MyMarketList'
import TablePagination from '@/components/TablePagination'
import DashTabs from '@/components/DashTabs'
import DashTab from '@/components/DashTab'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(0),
    },
    marginTop: theme.spacing(2),
  },
}))

const marketFilter = {
  type: MARKET_TYPE_ASK,
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
  const [filter, setFilter] = React.useState(marketFilter)
  const [reloadFlag, setReloadFlag] = React.useState(false)

  const [tabValue, setTabValue] = React.useState('/')

  const handleChange = (e, v) => {
    setTabValue(v)
  }

  React.useEffect(() => {
    ;(async () => {
      setData({ ...data, loading: true, error: null })
      try {
        const res = await myMarketSearch(filter)
        setData({ ...data, loading: false, ...res })
      } catch (e) {
        setData({ ...data, loading: false, error: e.message })
      }
    })()
  }, [filter, reloadFlag])

  React.useEffect(() => {
    ;(async () => {
      const res = await myMarketSearch(filter)
      setTotal(res.total_count)
    })()
  }, [])

  const handleSearchInput = value => {
    setFilter({ ...filter, loading: true, page: 1, q: value })
  }
  const handlePageChange = (e, page) => {
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
          <DashTabs value={tabValue} onChange={handleChange}>
            <DashTab
              component={Link}
              href="/my-dashboard"
              disableUnderline
              value="/"
              label="Active Listings"
              badgeContent={total}
            />
            <DashTab
              component={Link}
              href="/my-dashboard#reserved"
              disableUnderline
              value="/reserved"
              label="Reserved"
              badgeContent={12}
            />
            <DashTab
              component={Link}
              href="/my-dashboard#delivered"
              disableUnderline
              value="/delivered"
              label="Delivered"
              badgeContent={1}
            />
            <DashTab
              component={Link}
              href="/my-dashboard#history"
              disableUnderline
              value="/history"
              label="History"
            />
          </DashTabs>

          <MyMarketList
            datatable={data}
            loading={data.loading}
            error={data.error}
            onSearchInput={handleSearchInput}
            onReload={handleReloadToggle}
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
