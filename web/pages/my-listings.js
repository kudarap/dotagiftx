import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import WarningAmber from '@mui/icons-material/WarningAmber'
import {
  MARKET_STATUS_LIVE,
  MARKET_STATUS_RESERVED,
  MARKET_STATUS_SOLD,
  MARKET_TYPE_ASK,
} from '@/constants/market'
import { statsMarketSummary, marketSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import MyMarketList from '@/components/MyMarketList'
import DashTabs from '@/components/DashTabs'
import DashTab from '@/components/DashTab'
import AppContext from '@/components/AppContext'
import ReservationList from '@/components/ReservationList'
import MyMarketActivity from '@/components/MyMarketActivity'
import FixMyMarketActivity from '@/components/FixMyMarketActivity'
import withDatatableFetch from '@/components/withDatatableFetch'
import TabPanel from '@/components/TabPanel'

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

const marketIssueFilter = {
  partner_steam_id: 'https:',
}

export default function MyListings() {
  const { classes } = useStyles()

  const { currentAuth } = React.useContext(AppContext)

  // fetch market stats data
  const [marketStats, setMarketStats] = React.useState(initialMarketStats)
  const [tabValue, setTabValue] = React.useState(false)

  // fid market issues needs fixing
  const [marketIssues, setMarketFixes] = React.useState([])
  React.useEffect(() => {
    ;(async () => {
      const res = await marketSearch({
        user_id: currentAuth.user_id,
        index: 'user_id',
        ...marketIssueFilter,
      })
      setMarketFixes(res?.data || [])
    })()
  }, [currentAuth.user_id])

  // tick indicates when to get new stats
  const [tick, setTick] = React.useState(false)
  React.useEffect(() => {
    ;(async () => {
      const res = await statsMarketSummary({ user_id: currentAuth.user_id, index: 'user_id' })
      setMarketStats(res)
    })()
  }, [tick, currentAuth.user_id])

  // handling tab changes
  const router = useRouter()
  React.useEffect(() => {
    const hash = router.asPath.replace(router.pathname, '')
    setTabValue(hash)
  }, [router.asPath, router.pathname])

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
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            stats={marketStats}
            issues={marketIssues.length}
          />

          <TabPanel value={tabValue} index="">
            <LiveTable onReload={handleTableChange} />
          </TabPanel>
          <TabPanel value={tabValue} index="#reserved">
            <ReservedTable onReload={handleTableChange} />
          </TabPanel>
          <TabPanel value={tabValue} index="#delivered">
            <DeliveredTable />
          </TabPanel>
          <TabPanel value={tabValue} index="#history">
            <HistoryTable />
          </TabPanel>

          {marketIssues.length != 0 && (
            <TabPanel value={tabValue} index="#fix-broken-partner-steam-id">
              <BrokenPartnerSteamID onReload={handleTableChange} />
            </TabPanel>
          )}
        </Container>
      </main>

      <Footer />
    </>
  )
}

function Tabs(props) {
  const { stats, issues, ...other } = props

  return (
    <DashTabs {...other}>
      <DashTab value="" label="Active Listings" badgeContent={stats.live} />
      <DashTab value="#reserved" label="Reserved" badgeContent={stats.reserved} />
      <DashTab value="#delivered" label="Delivered" badgeContent={stats.sold} />
      <DashTab value="#history" label="History" />

      {issues != 0 && (
        <DashTab
          value="#fix-broken-partner-steam-id"
          label={
            <>
              <WarningAmber color="warning" fontSize="inherit" sx={{ mr: 0.6, mt: 0.25 }} />
              Fix Steam ID
            </>
          }
          badgeContent={issues}
        />
      )}
    </DashTabs>
  )
}
Tabs.propTypes = {
  stats: PropTypes.object.isRequired,
  issues: PropTypes.number.isRequired,
}

const datatableBaseFilter = {
  type: MARKET_TYPE_ASK,
  index: 'user_id',
}

const LiveTable = withDatatableFetch(MyMarketList, {
  ...datatableBaseFilter,
  status: MARKET_STATUS_LIVE,
})

const ReservedTable = withDatatableFetch(ReservationList, {
  ...datatableBaseFilter,
  status: MARKET_STATUS_RESERVED,
})

const DeliveredTable = withDatatableFetch(MyMarketActivity, {
  ...datatableBaseFilter,
  status: MARKET_STATUS_SOLD,
})

const HistoryTable = withDatatableFetch(MyMarketActivity, {
  ...datatableBaseFilter,
  limit: 20,
})

const BrokenPartnerSteamID = withDatatableFetch(FixMyMarketActivity, {
  ...datatableBaseFilter,
  ...marketIssueFilter,
  limit: 20,
})
