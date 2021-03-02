import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import useSWR from 'swr'
import has from 'lodash/has'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { fetcher, MARKETS, statsMarketSummary } from '@/service/api'
import MarketActivity from '@/components/MarketActivityV2'
import {
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_RESERVED,
  MARKET_STATUS_SOLD,
  MARKET_TYPE_ASK,
} from '@/constants/market'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  profile: {
    float: 'left',
    marginRight: theme.spacing(1),
    width: 60,
    height: 60,
  },
  itemImage: { width: 60, height: 40, marginRight: 8, float: 'left' },
}))

const filter = {
  type: MARKET_TYPE_ASK,
  sort: 'updated_at:desc',
  limit: 100,
}

export default function History({ status, summary }) {
  const classes = useStyles()

  filter.status = status
  const { data, error, isValidating } = useSWR([MARKETS, filter], fetcher, {
    revalidateOnFocus: false,
  })

  return (
    <>
      <Header />

      <Head>
        <title>
          {APP_NAME} :: Market {MARKET_STATUS_MAP_TEXT[status]} History
        </title>
        <meta name="description" content="Market transaction history" />
      </Head>

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1">
            Market History
          </Typography>

          <Typography style={{ display: 'flex' }}>
            <Typography component={Link} href="/history?reserved">
              {summary.reserved} Reserved
            </Typography>
            &nbsp;&middot;&nbsp;
            <Typography component={Link} href="/history?delivered">
              {summary.sold} Delivered
            </Typography>
          </Typography>

          {error && <Typography color="error">{error.message.split(':')[0]}</Typography>}
          <MarketActivity datatable={data || {}} loading={isValidating} disablePrice />
        </Container>
      </main>

      <Footer />
    </>
  )
}
History.propTypes = {
  status: PropTypes.number.isRequired,
  summary: PropTypes.object.isRequired,
}

export async function getServerSideProps({ query }) {
  let status = null
  if (has(query, 'reserved')) {
    status = MARKET_STATUS_RESERVED
  } else if (has(query, 'delivered')) {
    status = MARKET_STATUS_SOLD
  }

  const summary = await statsMarketSummary()

  return {
    props: { status, summary },
  }
}
