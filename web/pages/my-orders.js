import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import has from 'lodash/has'
import { makeStyles } from '@material-ui/core/styles'
import {
  MARKET_STATUS_BID_COMPLETED,
  MARKET_STATUS_LIVE,
  MARKET_STATUS_RESERVED,
  MARKET_TYPE_ASK,
  MARKET_TYPE_BID,
} from '@/constants/market'
import { marketSearch, statsMarketSummary } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import DashTabs from '@/components/DashTabs'
import DashTab from '@/components/DashTab'
import ReservationList from '@/components/ReservationList'
import MyMarketActivity from '@/components/MyMarketActivity'
import withDatatableFetch from '@/components/withDatatableFetch'
import AppContext from '@/components/AppContext'
import MyBuyOrderList from '@/components/MyBuyOrderList'

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

  const { currentAuth } = React.useContext(AppContext)
  console.log(currentAuth.steam_id)

  // fetch market stats data
  const router = useRouter()
  const [marketStats, setMarketStats] = React.useState(initialMarketStats)
  const [tabValue, setTabValue] = React.useState(false)
  React.useEffect(() => {
    ;(async () => {
      const res = await statsMarketSummary({ user_id: currentAuth.user_id })
      // Fetches count of linked markets.
      const linkedMarket = await statsMarketSummary({
        partner_steam_id: currentAuth.steam_id,
      })
      res.bids.reserved = linkedMarket.reserved
      setMarketStats(res.bids)
    })()
  }, [])

  // handling tab changes
  React.useEffect(() => {
    const hash = router.asPath.replace(router.pathname, '')
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
            <BuyOrdersTable />
          </TabPanel>
          <TabPanel value={tabValue} index="#toreceive">
            <ToReceiveTable filter={{ partner_steam_id: currentAuth.steam_id }} />
          </TabPanel>
          <TabPanel value={tabValue} index="#completed">
            <CompletedTable />
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
      <DashTab value="" label="Buy Orders" badgeContent={stats.live} />
      <DashTab value="#toreceive" label="To Receive" badgeContent={stats.reserved} />
      <DashTab value="#completed" label="Completed" badgeContent={stats.bid_completed} />
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

const datatableBaseFilter = {
  type: MARKET_TYPE_BID,
}

const BuyOrdersTable = withDatatableFetch(MyBuyOrderList, {
  ...datatableBaseFilter,
  status: MARKET_STATUS_LIVE,
})
const ToReceiveTable = withDatatableFetch(
  ReservationList,
  {
    type: MARKET_TYPE_ASK,
    status: MARKET_STATUS_RESERVED,
  },
  marketSearch
)
const CompletedTable = withDatatableFetch(MyMarketActivity, {
  ...datatableBaseFilter,
  status: MARKET_STATUS_BID_COMPLETED,
})
const HistoryTable = withDatatableFetch(MyMarketActivity, {
  ...datatableBaseFilter,
  limit: 20,
})