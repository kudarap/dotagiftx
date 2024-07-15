import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import Head from 'next/head'
import has from 'lodash/has'
import { makeStyles } from 'tss-react/mui'
import Avatar from '@/components/Avatar'
import Typography from '@mui/material/Typography'
import {
  APP_NAME,
  APP_URL,
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import { USER_STATUS_MAP_TEXT } from '@/constants/user'
import { MARKET_STATUS_LIVE, MARKET_TYPE_ASK } from '@/constants/market'
import {
  CDN_URL,
  isDonationGlowExpired,
  marketSearch,
  trackProfileViewURL,
  user,
  vanity,
} from '@/service/api'
import { getUserBadgeFromBoons } from '@/lib/badge'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Container from '@/components/Container'
import ChipLink from '@/components/ChipLink'
import UserMarketList from '@/components/UserMarketList'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import Link from '@/components/Link'
import Button from '@/components/Button'
import NotRegisteredProfile from '@/components/NotRegisteredProfile'
import AppContext from '@/components/AppContext'
import SubscriberBadge from '@/components/SubscriberBadge'
import ErrorPage from '../404'
import { Alert } from '@mui/material'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  profileName: {
    [theme.breakpoints.down('sm')]: {
      fontSize: theme.typography.h6.fontSize,
    },
  },
  details: {
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  avatar: {
    [theme.breakpoints.down('sm')]: {
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
  const { classes } = useStyles()

  const [markets, setMarkets] = React.useState(initialMarkets)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState(initialError)
  const { isMobile } = React.useContext(AppContext)

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

  const userBadge = getUserBadgeFromBoons(profile.boons)

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
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
          {isProfileReported && (
            <>
              <Alert
                severity="error"
                variant="outlined"
                sx={{
                  fontSize: '1rem',
                  borderColor: '#c13830',
                  borderWidth: 2,
                }}>
                This is user has been flagged as <strong>BANNED</strong> or{' '}
                <strong>SUSPENDED</strong>. <br />
                Website is not liable for any lost in-game items and money and should avoid any
                transaction.
              </Alert>
              <br />
            </>
          )}

          <div
            className={classes.details}
            style={
              isProfileReported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : null
            }>
            <Avatar
              large
              badge={userBadge}
              className={classes.avatar}
              src={`${CDN_URL}/${profile.avatar}`}
              glow={isDonationGlowExpired(profile.donated_at)}
            />
            <Typography component="h1">
              <Typography
                className={classes.profileName}
                component="p"
                variant="h4"
                color={isProfileReported ? 'error' : 'textPrimary'}>
                {profile.name}
                {Boolean(userBadge) && (
                  <SubscriberBadge
                    style={
                      isMobile
                        ? { margin: '0 4px' }
                        : { marginLeft: 4, marginTop: 12, position: 'absolute' }
                    }
                    type={userBadge}
                    size="medium"
                  />
                )}
              </Typography>
              {isProfileReported && (
                <Typography color="error">
                  {profile.notes || USER_STATUS_MAP_TEXT[profile.status]}
                </Typography>
              )}
              <Typography gutterBottom>
                <Typography variant="body2" component="span">
                  <Link href={`${linkProps.href}`}>{profile.stats.live} Items</Link> &middot;{' '}
                  <Link href={`${linkProps.href}/reserved`}>{profile.stats.reserved} Reserved</Link>{' '}
                  &middot;{' '}
                  <Link href={`${linkProps.href}/delivered`}>{profile.stats.sold} Delivered</Link>{' '}
                  &middot;{' '}
                  <Link href={`${linkProps.href}/bought`}>
                    {profile.stats.bid_completed} Bought
                  </Link>
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

        <img src={trackProfileViewURL(profile.id)} height={1} width={1} alt="" />
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
  index: 'user_id',
}

// This gets called on every request
export async function getServerSideProps(context) {
  const { params, query } = context
  const vanityMode = Boolean(query.vanity)

  let profile
  let canonicalURL

  // Check for vanity request.
  if (vanityMode) {
    try {
      profile = await vanity(String(query.vanity))
      canonicalURL = `${APP_URL}/id/${query.vanity}`
      return {
        redirect: {
          permanent: true,
          destination: `/profiles/${profile.steam_id}`,
        },
      }
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
  profile.stats = profile.market_stats

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
