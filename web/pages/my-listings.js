import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import { myMarketSearch, statsMarketSummary } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import {
  MARKET_STATUS_LIVE,
  MARKET_STATUS_RESERVED,
  MARKET_STATUS_SOLD,
  MARKET_TYPE_ASK,
} from '@/constants/market'
import MyMarketList from '@/components/MyMarketList'
import TablePagination from '@/components/TablePagination'
import DashTabs from '@/components/DashTabs'
import DashTab from '@/components/DashTab'
import { useRouter } from 'next/router'
import AppContext from '@/components/AppContext'
import ReservationList from '@/components/ReservationList'
import HistoryList from '@/components/HistoryList'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(2),
  },
}))

const initialMarketStats = {
  pending: 0,
  live: 0,
  reserved: 0,
  sold: 0,
}

export default function MyListings() {
  const classes = useStyles()

  // fetch market stats data
  const router = useRouter()
  const { currentAuth } = React.useContext(AppContext)
  const [marketStats, setMarketStats] = React.useState(initialMarketStats)
  const [tabValue, setTabValue] = React.useState('')

  React.useEffect(() => {
    ;(async () => {
      const res = await statsMarketSummary({ user_id: currentAuth.user_id })
      setMarketStats(res)
    })()
  }, [currentAuth])

  // handling tab changes
  React.useEffect(() => {
    setTabValue(router.asPath.replace(router.pathname, ''))
  }, [router.asPath])

  const handleTabChange = (e, v) => {
    setTabValue(v)
    router.push(v)
  }

  // const tabContents = {
  //   '': <LiveTable />,
  //   '#reserved': <ReservedTable />,
  //   '#delivered': <DeliveredTable />,
  //   '#history': <HistoryTable />,
  // }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Tabs value={tabValue} onChange={handleTabChange} stats={marketStats} />
          <LiveTable />
        </Container>
      </main>

      <Footer />
    </>
  )
}

function Tabs(props) {
  const { stats, ...other } = props

  return (
    <DashTabs {...other}>
      <DashTab value="" label="Active Listings" badgeContent={stats.live} />
      <DashTab value="#reserved" label="Reserved" badgeContent={stats.reserved} />
      <DashTab value="#delivered" label="Delivered" badgeContent={stats.sold} />
      <DashTab value="#history" label="History" />
    </DashTabs>
  )
}

const initialDatatable = {
  data: [],
  result_count: 0,
  total_count: 0,
  loading: false,
  error: null,
}

const datatableBaseFilter = {
  type: MARKET_TYPE_ASK,
  sort: 'created_at:desc',
  page: 1,
}

const withDataFetch = (Component, initialFilter) => props => {
  const [data, setData] = React.useState(initialDatatable)
  const [filter, setFilter] = React.useState({ ...datatableBaseFilter, ...initialFilter })
  const [tick, setTick] = React.useState(false)

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
  }, [filter, tick])

  const handleSearchInput = value => {
    setFilter({ ...filter, loading: true, page: 1, q: value })
  }
  const handlePageChange = (e, page) => {
    setFilter({ ...filter, page })
  }
  const handleReloadToggle = () => {
    setTick(!tick)
  }

  return (
    <>
      <Component
        datatable={data}
        loading={data.loading}
        error={data.error}
        onSearchInput={handleSearchInput}
        onReload={handleReloadToggle}
        {...props}
      />
      <TablePagination
        style={{ textAlign: 'right' }}
        count={data.total_count || 0}
        page={filter.page}
        onChangePage={handlePageChange}
      />
    </>
  )
}

const LiveTable = withDataFetch(MyMarketList, { status: MARKET_STATUS_LIVE })
const ReservedTable = withDataFetch(ReservationList, { status: MARKET_STATUS_RESERVED })
const DeliveredTable = withDataFetch(HistoryList, { status: MARKET_STATUS_SOLD })
const HistoryTable = withDataFetch(HistoryList, {})
