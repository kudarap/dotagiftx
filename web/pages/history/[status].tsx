import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import * as format from '@/lib/format'
import { marketSearch, statsMarketSummaryOverall } from '@/service/api'
import { APP_NAME } from '@/constants/strings'
import {
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_RESERVED,
  MARKET_STATUS_SOLD,
  MARKET_TYPE_ASK,
} from '@/constants/market'
import Link from '@/components/Link'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import MarketActivity from '@/components/MarketActivity'
import type { MarketSummary } from '@/lib/types'
import type { ActivityMarket } from '@/components/MarketActivity'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
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
  itemImage: {
    width: 60,
    height: 40,
    marginRight: 8,
    float: 'left',
  },
  nav: {
    display: 'flex',
    '& active[]': {
      color: 'white',
    },
    marginBottom: theme.spacing(2),
  },
}))

const defaultFilter = {
  type: MARKET_TYPE_ASK,
  sort: 'updated_at:desc',
  page: 1,
  limit: 15,
}

const defaultData = {
  data: [] as ActivityMarket[],
  total_count: 0,
}

const scrollBias = 300

export default function History({
  status,
  summary,
  error,
}: {
  status: number | null
  summary: MarketSummary
  error?: string | null
}) {
  const { classes } = useStyles()

  const [datatable, setDatatable] = React.useState(defaultData)
  const [filter, setFilter] = React.useState({ ...defaultFilter, status })
  const [loading, setLoading] = React.useState(false)
  const [filterError, setFilterError] = React.useState('')

  const busyRef = React.useRef(false)

  React.useEffect(() => {
    setDatatable(defaultData)
    setFilter({ ...defaultFilter, status })
  }, [status])

  React.useEffect(() => {
    if (busyRef.current) {
      return
    }

    busyRef.current = true
    setLoading(true)
    ;(async () => {
      try {
        const res = (await marketSearch(filter)) as { data: ActivityMarket[]; total_count: number }
        setDatatable(current =>
          current.data.length === 0 ? res : { ...current, data: [...current.data, ...res.data] }
        )
      } catch (e) {
        setFilterError(`error getting history ${(e as Error).message}`)
      }

      busyRef.current = false
      setLoading(false)
    })()
  }, [filter])

  React.useEffect(() => {
    const listener = () => {
      const isLast = datatable.data.length === datatable.total_count
      const scrollMaxY = (window as unknown as { scrollMaxY: number }).scrollMaxY
      if (loading || isLast || window.scrollY + scrollBias < scrollMaxY) {
        return
      }

      setFilter({ ...filter, page: (filter.page as number) + 1 })
    }

    window.addEventListener('scroll', listener)
    return () => {
      window.removeEventListener('scroll', listener)
    }
  })

  const summarySold = format.numberWithCommas(summary.sold)
  const summaryReserved = format.numberWithCommas(summary.reserved)
  return (
    <>
      <Header />

      <Head>
        <meta charSet="UTF-8" />
        <title>
          {APP_NAME} :: Market {status !== null ? MARKET_STATUS_MAP_TEXT[status as number] : ''}{' '}
          History
        </title>
        <meta name="description" content="Market transaction history" />
      </Head>

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1">
            Market History
          </Typography>

          <Typography className={classes.nav}>
            <Typography
              component={Link}
              href="/history/reserved"
              style={
                status === MARKET_STATUS_RESERVED ? { textDecoration: 'underline' } : undefined
              }
            >
              {summaryReserved} Reserved
            </Typography>
            &nbsp;&middot;&nbsp;
            <Typography
              component={Link}
              href="/history/delivered"
              style={status === MARKET_STATUS_SOLD ? { textDecoration: 'underline' } : undefined}
            >
              {summarySold} Delivered
            </Typography>
          </Typography>

          {filterError && (
            <Typography align="center" variant="body2" color="error">
              {filterError}
            </Typography>
          )}

          {error && <Typography color="error">{error.split(':')[0]}</Typography>}
          <MarketActivity datatable={datatable} loading={loading} disablePrice={status !== null} />
        </Container>
      </main>

      <Footer />
    </>
  )
}

export async function getServerSideProps({ query }: { query: { status?: string } }) {
  let status: number | null = null

  switch (query.status) {
    case 'reserved':
      status = MARKET_STATUS_RESERVED
      break
    case 'delivered':
      status = MARKET_STATUS_SOLD
      break
  }

  let summary: MarketSummary | null = null
  let error: string | null = null
  try {
    summary = (await statsMarketSummaryOverall()) as MarketSummary
  } catch (e) {
    error = (e as Error).message
  }

  return {
    props: { status, summary, error },
  }
}
