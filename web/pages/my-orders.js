import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
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
import MyMarketActivity from '@/components/MyMarketActivity'
import withDatatableFetch from '@/components/withDatatableFetch'
import AppContext from '@/components/AppContext'
import MyBuyOrderList from '@/components/MyBuyOrderList'
import TabPanel from '@/components/TabPanel'
import MyOrderActivity from '@/components/MyOrderActivity'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
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

export default function MyOrders() {
  const { classes } = useStyles()

  const { currentAuth } = React.useContext(AppContext)

  // fetch market stats data
  const [marketStats, setMarketStats] = React.useState(initialMarketStats)
  const [tabValue, setTabValue] = React.useState(false)

  // tick indicates when to get new stats
  const [tick, setTick] = React.useState(false)

  React.useEffect(() => {
    ;(async () => {
      const res = await statsMarketSummary({ user_id: currentAuth.user_id })
      // Fetches count of linked markets.
      const linkedMarket = await statsMarketSummary({ partner_steam_id: currentAuth.steam_id })
      res.bids.reserved = linkedMarket.reserved
      setMarketStats(res.bids)
    })()
  }, [tick])

  // handling tab changes
  const router = useRouter()
  React.useEffect(() => {
    const hash = router.asPath.replace(router.pathname, '')
    setTabValue(hash)
  }, [router.asPath])

  const handleTabChange = (e, v) => {
    setTabValue(v)
    router.push(v)
  }

  const handleTableChange = () => {
    setTick(!tick)
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Tabs value={tabValue} onChange={handleTabChange} stats={marketStats} />

          <TabPanel value={tabValue} index="">
            <BuyOrdersTable onReload={handleTableChange} />
          </TabPanel>
          <TabPanel value={tabValue} index="#toreceive">
            <ToReceiveTable
              filter={{
                index: 'partner_steam_id',
                partner_steam_id: currentAuth.steam_id,
              }}
            />
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

const datatableBaseFilter = {
  type: MARKET_TYPE_BID,
}

const BuyOrdersTable = withDatatableFetch(MyBuyOrderList, {
  ...datatableBaseFilter,
  status: MARKET_STATUS_LIVE,
})
const ToReceiveTable = withDatatableFetch(
  MyOrderActivity,
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
