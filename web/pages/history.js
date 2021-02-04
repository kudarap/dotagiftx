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
import { fetcher, MARKETS } from '@/service/api'
import MarketActivity from '@/components/MarketActivity'
import {
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_RESERVED,
  MARKET_STATUS_SOLD,
  MARKET_TYPE_ASK,
} from '@/constants/market'

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

export default function History({ status }) {
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
          <Typography component="h1">
            {MARKET_STATUS_MAP_TEXT[status]} Items {data && `(${data && data.total_count})`}
            {/* <CircularProgress color="secondary" size={15} /> */}
          </Typography>
          {error && <Typography color="error">{error.message.split(':')[0]}</Typography>}
          <MarketActivity data={data ? data.data : null} loading={isValidating} />
        </Container>
      </main>

      <Footer />
    </>
  )
}
History.propTypes = {
  status: PropTypes.number.isRequired,
}

export async function getServerSideProps({ query }) {
  let status = null
  if (has(query, 'reserved')) {
    status = MARKET_STATUS_RESERVED
  } else if (has(query, 'delivered')) {
    status = MARKET_STATUS_SOLD
  }

  return {
    props: { status },
  }
}
