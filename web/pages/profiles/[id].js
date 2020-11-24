import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { CDN_URL, marketSearch, statsMarketSummary, user } from '@/service/api'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Container from '@/components/Container'
import ChipLink from '@/components/ChipLink'
import UserMarketList from '@/components/UserMarketList'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import {
  APP_NAME,
  APP_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  profileName: {
    [theme.breakpoints.down('xs')]: {
      fontSize: theme.typography.h6.fontSize,
    },
  },
  details: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  avatar: {
    [theme.breakpoints.down('xs')]: {
      margin: '0 auto',
    },
    width: 100,
    height: 100,
    marginRight: theme.spacing(1.5),
  },
}))

export default function UserDetails({
  profile,
  filter,
  markets: initialMarkets,
  error: initialError,
  canonicalURL,
}) {
  const classes = useStyles()

  const [markets, setMarkets] = React.useState(initialMarkets)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState(initialError)

  // Handle market request on page change.
  React.useEffect(() => {
    ;(async () => {
      setError(null)
      setLoading(true)
      try {
        const res = await marketSearch(filter)
        setMarkets(res)
      } catch (e) {
        setError(e.message)
      }
      setLoading(false)
    })()
  }, [filter])

  const router = useRouter()
  const qFilter = router.query.filter
  const linkProps = {
    href: `/profiles/${profile.steam_id}`,
  }
  if (String(qFilter).trim() !== '') {
    linkProps.query = { filter: qFilter }
  }

  const handleSearchInput = text => {
    let url = `${linkProps.href}?filter=${text}`
    if (String(text).trim() === '') {
      url = linkProps.href
    }

    router.push(url)
  }

  const profileURL = `${STEAM_PROFILE_BASE_URL}/${profile.steam_id}`
  const steamRepURL = `${STEAMREP_PROFILE_BASE_URL}/${profile.steam_id}`

  const metaTitle = `${APP_NAME} :: ${profile.name}`
  const metaDesc = `${profile.name}'s Dota 2 giftable item listings`

  return (
    <>
      <Head>
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={canonicalURL} />

        {/* Twitter Card */}
        <meta name="twitter:card" content="summary" />
        <meta name="twitter:title" content={metaTitle} />
        <meta name="twitter:description" content={metaDesc} />
        <meta name="twitter:image" content={`${CDN_URL}/${profile.avatar}`} />
        <meta name="twitter:site" content={`@${APP_NAME}`} />
        {/* OpenGraph */}
        <meta property="og:url" content={canonicalURL} />
        <meta property="og:type" content="website" />
        <meta property="og:title" content={metaTitle} />
        <meta property="og:description" content={metaDesc} />
        <meta property="og:image" content={`${CDN_URL}/${profile.avatar}`} />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <div className={classes.details}>
            <Avatar className={classes.avatar} src={`${CDN_URL}/${profile.avatar}`} />
            <Typography component="h1">
              <Typography className={classes.profileName} component="p" variant="h4">
                {profile.name}
              </Typography>
              <Typography gutterBottom>
                <Typography variant="body2" component="span">
                  <Link href={`${linkProps.href}`}>{profile.stats.live} Items</Link> &middot;{' '}
                  <Link href={`${linkProps.href}/reserved`}>{profile.stats.reserved} Reserved</Link>{' '}
                  &middot;{' '}
                  <Link href={`${linkProps.href}/delivered`}>{profile.stats.sold} Delivered</Link>
                </Typography>
                <br />
                <ChipLink label="Steam Profile" href={profileURL} />
                &nbsp;
                {/* <ChipLink label="Steam Inventory" href={`${profileURL}/inventory`} /> */}
                {/* &nbsp; */}
                <ChipLink label="SteamRep" href={steamRepURL} />
              </Typography>
            </Typography>
          </div>

          <UserMarketList
            onSearchInput={handleSearchInput}
            data={markets}
            loading={loading}
            error={error}
          />
          {!error && (
            <TablePaginationRouter
              linkProps={linkProps}
              style={{ textAlign: 'right' }}
              count={markets.total_count}
              page={filter.page}
            />
          )}
        </Container>
      </main>

      <Footer />
    </>
  )
}
UserDetails.propTypes = {
  profile: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
  filter: PropTypes.object,
  markets: PropTypes.object,
  error: PropTypes.string,
}
UserDetails.defaultProps = {
  filter: {},
  markets: {
    data: [],
  },
  error: null,
}

const marketSearchFilter = {
  page: 1,
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
}

// This gets called on every request
export async function getServerSideProps({ params, query }) {
  const profile = await user(String(params.id))
  const filter = { ...marketSearchFilter, user_id: profile.id }
  filter.page = Number(query.page || 1)
  if (query.filter) {
    filter.q = query.filter
  }

  profile.stats = await statsMarketSummary({ user_id: profile.id })

  let markets = {}
  let error = null
  try {
    markets = await marketSearch(filter)
  } catch (e) {
    error = e.message
  }

  const canonicalURL = `${APP_URL}/profiles/${params.id}`

  return {
    props: {
      profile,
      canonicalURL,
      filter,
      markets,
      error,
    },
  }
}
