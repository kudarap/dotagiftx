import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import Head from 'next/head'
import has from 'lodash/has'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@/components/Avatar'
import Typography from '@material-ui/core/Typography'
import DonatorIcon from '@material-ui/icons/FavoriteBorder'
import { MARKET_STATUS_LIVE, MARKET_TYPE_ASK } from '@/constants/market'
import {
  CDN_URL,
  marketSearch,
  statsMarketSummary,
  trackProfileViewURL,
  user,
  vanity,
} from '@/service/api'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Container from '@/components/Container'
import ChipLink from '@/components/ChipLink'
import UserMarketList from '@/components/UserMarketList'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import {
  APP_NAME,
  APP_URL,
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import { USER_STATUS_MAP_TEXT } from '@/constants/user'
import Link from '@/components/Link'
import Button from '@/components/Button'
import NotRegisteredProfile from '@/components/NotRegisteredProfile'
import ErrorPage from '../404'
import DonatorBadge from '@/components/DonatorBadge'

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
    marginBottom: theme.spacing(0.5),
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

  if (error) {
    return (
      <ErrorPage>
        <Typography variant="h5" align="center">
          Profile not found
        </Typography>
      </ErrorPage>
    )
  }

  // This user is not registered
  if (has(profile, 'is_registered') && !profile.is_registered) {
    return <NotRegisteredProfile profile={profile} canonicalURL={canonicalURL} />
  }

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
  const dotabuffURL = `${DOTABUFF_PROFILE_BASE_URL}/${profile.steam_id}`

  const metaTitle = `${APP_NAME} :: ${profile.name}`
  let metaDesc = `${profile.name}'s Dota 2 Giftable`
  if (profile.stats) {
    metaDesc += ` ${profile.stats.live} Items · ${profile.stats.reserved} Reserved · ${profile.stats.sold} Delivered`
  }

  const isProfileReported = Boolean(profile.status)

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
          <div
            className={classes.details}
            style={
              isProfileReported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : null
            }>
            <Avatar
              className={classes.avatar}
              src={`${CDN_URL}/${profile.avatar}`}
              glow={Boolean(profile.donation)}
            />
            <Typography component="h1">
              <Typography
                className={classes.profileName}
                component="p"
                variant="h4"
                color={isProfileReported ? 'error' : 'textPrimary'}>
                {profile.name}
                {Boolean(profile.donation) && (
                  <DonatorBadge
                    style={{ marginLeft: 4, marginTop: 12, position: 'absolute' }}
                    size="medium">
                    DONATOR
                  </DonatorBadge>
                )}
              </Typography>
              {isProfileReported && (
                <Typography color="error">{USER_STATUS_MAP_TEXT[profile.status]}</Typography>
              )}
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
                &nbsp;
                <ChipLink label="Dotabuff" href={dotabuffURL} />
              </Typography>
            </Typography>
          </div>

          {isProfileReported ? (
            <p align="center">
              <Button component={Link} href={`${linkProps.href}/activity`}>
                Show All Activity
              </Button>
            </p>
          ) : (
            <>
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
            </>
          )}
        </Container>

        <img src={trackProfileViewURL(profile.id)} alt="" />
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
  type: MARKET_TYPE_ASK,
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
}

// This gets called on every request
export async function getServerSideProps({ params, query }) {
  const vanityMode = Boolean(query.vanity)

  let profile
  let canonicalURL

  // Check for vanity request.
  if (vanityMode) {
    try {
      profile = await vanity(String(query.vanity))
      canonicalURL = `${APP_URL}/id/${query.vanity}`
    } catch (e) {
      return {
        props: {
          error: e.message,
        },
      }
    }

    // Since not registered user will render differently, should return now.
    if (!profile.is_registered) {
      return {
        props: {
          profile,
          canonicalURL,
        },
      }
    }

    // When vanity exists use the profile from resolving it.
    // Otherwise try to get from users endpoint
  } else {
    try {
      profile = await user(String(params.id))
    } catch (e) {
      return {
        props: {
          error: e.message,
        },
      }
    }
  }

  // Retrieve initial user market summary.
  profile.stats = await statsMarketSummary({ type: MARKET_TYPE_ASK, user_id: profile.id })

  // Retrieve initial user market data.
  let markets = {}
  let error = null
  const filter = { ...marketSearchFilter, user_id: profile.id }
  filter.page = Number(query.page || 1)
  if (query.filter) {
    filter.q = query.filter
  }

  try {
    markets = await marketSearch(filter)
  } catch (e) {
    error = e.message
  }

  // Compose profile page canonical URL.

  canonicalURL = `${APP_URL}/${vanityMode ? 'id' : 'profiles'}/${params.id}`

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
