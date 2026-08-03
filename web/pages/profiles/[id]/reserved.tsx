import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Avatar from '@/components/Avatar'
import { APP_NAME, APP_URL } from '@/constants/strings'
import { MARKET_STATUS_RESERVED } from '@/constants/market'
import { CDN_URL, marketSearch, user } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import MarketActivity from '@/components/MarketActivity'
import type { Profile } from '@/lib/types'
import type { ActivityMarket } from '@/components/MarketActivity'
import type { QueryFilter } from '@/service/http'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(1),
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

const defaultFilter = {
  status: MARKET_STATUS_RESERVED,
  sort: 'updated_at:desc',
  page: 1,
  index: 'user_id',
}

const defaultData = {
  data: [] as ActivityMarket[],
  total_count: 0,
}

const scrollBias = 300

export default function UserReserved({
  profile,
  stats,
  canonicalURL,
}: {
  profile: Profile
  stats: Profile['stats']
  canonicalURL: string
}) {
  const { classes } = useStyles()

  const [datatable, setDatatable] = React.useState(defaultData)
  const [filter, setFilter] = React.useState<QueryFilter & { page?: number }>({
    ...defaultFilter,
    user_id: profile.id,
  })
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<Error | null>(null)

  const busyRef = React.useRef(false)

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
        setError(e as Error)
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

      setFilter({ ...filter, page: (filter.page || 1) + 1 })
    }

    window.addEventListener('scroll', listener)
    return () => {
      window.removeEventListener('scroll', listener)
    }
  })

  const handleSearchInput = (q: string) => {
    setDatatable(defaultData)
    setFilter({ ...filter, page: 1, q })
  }

  const profileURL = `/profiles/${profile.steam_id}`

  const s = stats || {
    live: 0,
    reserved: 0,
    sold: 0,
    bid_completed: 0,
  }

  return (
    <>
      <Header />

      <Head>
        <meta charSet="UTF-8" />
        <title>{`${APP_NAME} :: ${profile.name} reserved items`}</title>
        <meta name="description" content={`${profile.name}'s on-reserved Giftable items`} />
        <link rel="canonical" href={canonicalURL} />
      </Head>

      <main className={classes.main}>
        <Container>
          <div>
            <Avatar
              className={classes.profile}
              src={`${CDN_URL}/${profile.avatar}`}
              glow={Boolean(profile.donation)}
              component={Link}
              href={profileURL}
            />
            <Typography variant="h6" color="textPrimary" component={Link} href={profileURL}>
              {profile.name}
            </Typography>
            <div style={{ display: 'flex' }}>
              <Typography component={Link} href={profileURL}>
                {s.live} Items
              </Typography>
              &nbsp;&middot;&nbsp;
              <Typography
                component={Link}
                href={`${profileURL}/reserved`}
                style={{ textDecoration: 'underline' }}
              >
                {s.reserved} Reserved
              </Typography>
              &nbsp;&middot;&nbsp;
              <Typography component={Link} href={`${profileURL}/delivered`}>
                {s.sold} Delivered
              </Typography>
              &nbsp;&middot;&nbsp;
              <Typography component={Link} href={`${profileURL}/bought`}>
                {s.bid_completed} Bought
              </Typography>
            </div>
          </div>
          <br />

          {error && <Typography color="error">{error.message.split(':')[0]}</Typography>}
          <MarketActivity
            datatable={datatable}
            loading={loading}
            onSearchInput={handleSearchInput}
            disablePrice
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}

export async function getServerSideProps({ params }: { params: { id: string } }) {
  const profile = (await user(String(params.id))) as Profile
  const canonicalURL = `${APP_URL}/profiles/${params.id}/reserved`

  const stats = profile.market_stats || {}

  return {
    props: {
      profile,
      canonicalURL,
      stats,
    },
  }
}
