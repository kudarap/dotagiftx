import React from 'react'
import PropTypes from 'prop-types'
import has from 'lodash/has'
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
import MyMarketActivity from '@/components/MyMarketActivity'

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
  const [tabValue, setTabValue] = React.useState(false)

  React.useEffect(() => {
    ;(async () => {
      const res = await statsMarketSummary({ user_id: currentAuth.user_id })
      setMarketStats(res)
    })()
  }, [])

  // handling tab changes
  React.useEffect(() => {
    const hash = router.asPath.replace(router.pathname, '')
    console.log(hash)
    setTabValue(hash)
  }, [router.asPath])

  const handleTabChange = (e, v) => {
    setTabValue(v)
    router.push(v)
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Tabs value={tabValue} onChange={handleTabChange} stats={marketStats} />

          <TabPanel value={tabValue} index="">
            <LiveTable />
          </TabPanel>
          <TabPanel value={tabValue} index="#reserved">
            <ReservedTable />
          </TabPanel>
          <TabPanel value={tabValue} index="#delivered">
            <DeliveredTable />
          </TabPanel>
          <TabPanel value={tabValue} index="#history">
            <HistoryTable />
          </TabPanel>
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
Tabs.propTypes = {
  stats: PropTypes.object.isRequired,
}

const tabPanelIndex = {}
function TabPanel(props) {
  const { children, value, index, ...other } = props

  // Check for indexed component, it will prevent render from
  // loading everything on mount.
  if (value !== index && !has(tabPanelIndex, index)) {
    return null
  }
  tabPanelIndex[index] = true

  return (
    <div hidden={value !== index} {...other}>
      {children}
    </div>
  )
}
TabPanel.propTypes = {
  children: PropTypes.node.isRequired,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
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

const withDataFetch = (Component, initFilter) => props => {
  const [data, setData] = React.useState(initialDatatable)
  const [filter, setFilter] = React.useState({ ...datatableBaseFilter, ...initFilter })
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
        rowsPerPage={filter.limit}
        onChangePage={handlePageChange}
      />
    </>
  )
}

const LiveTable = withDataFetch(MyMarketList, { status: MARKET_STATUS_LIVE })
const ReservedTable = withDataFetch(ReservationList, { status: MARKET_STATUS_RESERVED })
const DeliveredTable = withDataFetch(MyMarketActivity, { status: MARKET_STATUS_SOLD })
const HistoryTable = withDataFetch(MyMarketActivity, { limit: 20 })
