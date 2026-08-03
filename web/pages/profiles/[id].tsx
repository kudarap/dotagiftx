import React from 'react'
import { useRouter } from 'next/router'
import Head from 'next/head'
import has from 'lodash/has'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Alert from '@mui/material/Alert'
import moment from 'moment'
import NewbieIcon from '@mui/icons-material/NewReleases'
import { Box } from '@mui/material'
import {
  APP_NAME,
  APP_URL,
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
} from '@/constants/strings'
import {
  USER_AGE_CAUTION,
  USER_STATUS_MAP_TEXT,
  USER_SUBSCRIPTION_BADGE_MODE,
} from '@/constants/user'
import { MARKET_STATUS_LIVE, MARKET_TYPE_ASK } from '@/constants/market'
import {
  CDN_URL,
  isDonationGlowExpired,
  marketSearch,
  trackProfileViewURL,
  user,
  vanity,
} from '@/service/api'
import { getUserBadgeFromBoons, getUserTagFromBoons } from '@/lib/badge'
import Avatar from '@/components/Avatar'
import ExclusiveChip from '@/components/ExclusiveChip'
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
import type { Profile, Market } from '@/lib/types'
import type { MarketDatatable } from '@/components/MarketList'
import type { SubscriberBadgeType } from '@/components/SubscriberBadge'
import type { ExclusiveTagType } from '@/components/ExclusiveChip'

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
    width: 110,
    height: 110,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(0.5),
  },
}))

export default function UserDetails({
  profile,
  filter,
  markets: initialMarkets,
  profileError,
  marketError,
  canonicalURL,
}: {
  profile: Profile
  filter: { page?: number; q?: string; user_id?: string | number }
  markets: MarketDatatable
  profileError?: string | null
  marketError?: string | null
  canonicalURL: string
}) {
  const { classes } = useStyles()

  const [markets, setMarkets] = React.useState(initialMarkets)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(marketError ?? null)
  const { isMobile } = React.useContext(AppContext)

  const router = useRouter()

  // Handle market request on page change.
  React.useEffect(() => {
    if (profileError || (has(profile, 'is_registered') && !profile.is_registered)) {
      return
    }

    ;(async () => {
      setError(null)
      setLoading(true)
      try {
        const res = (await marketSearch(filter)) as MarketDatatable
        setMarkets(res)
      } catch (e) {
        setError((e as Error).message)
      }
      setLoading(false)
    })()
  }, [filter, profileError, profile])

  if (profileError) {
    return (
      <ErrorPage>
        <Typography variant="h5" align="center">
          {profileError}
        </Typography>
      </ErrorPage>
    )
  }

  // This user is not registered
  if (has(profile, 'is_registered') && !profile.is_registered) {
    const unregistered = profile as unknown as {
      steam_id: string
      name: string
      steam_avatar: string
      url: string
      last_updated_at: string
    }
    return <NotRegisteredProfile profile={unregistered} canonicalURL={canonicalURL} />
  }

  const qFilter = router.query.filter
  const linkProps: { href: string; query?: Record<string, string> } = {
    href: `/profiles/${profile.steam_id}`,
  }
  if (String(qFilter).trim() !== '') {
    linkProps.query = { filter: String(qFilter) }
  }

  const handleSearchInput = (text: string) => {
    let url = `${linkProps.href}?filter=${text}`
    if (String(text).trim() === '') {
      url = linkProps.href
    }

    router.push(url)
  }

  const profileURL = `${STEAM_PROFILE_BASE_URL}/${profile.steam_id}`
  const dotabuffURL = `${DOTABUFF_PROFILE_BASE_URL}/${profile.steam_id}`

  const metaTitle = `${APP_NAME} :: ${profile.name}`
  let metaDesc = `${profile.name}'s Dota 2 Giftable`
  if (profile.stats) {
    metaDesc += ` ${profile.stats.live} Items · ${profile.stats.reserved} Reserved · ${profile.stats.sold} Delivered`
  }

  const isProfileReported = Boolean(profile.status)

  const userBadge = getUserBadgeFromBoons(profile.boons) as SubscriberBadgeType | undefined
  const userTag = getUserTagFromBoons(profile.boons) as ExclusiveTagType | undefined

  const stats = profile.stats || {
    live: 0,
    reserved: 0,
    sold: 0,
    bid_completed: 0,
  }

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
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
        <meta property="og:site_name" content={APP_NAME} />
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
                }}
              >
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
              isProfileReported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : {}
            }
          >
            <Avatar
              large
              badge={userBadge}
              className={classes.avatar}
              src={`${CDN_URL}/${profile.avatar}`}
              glow={isDonationGlowExpired(profile.donated_at)}
            />
            <Box>
              <Typography
                className={classes.profileName}
                component="h1"
                variant="h4"
                color={isProfileReported ? 'error' : 'textPrimary'}
              >
                {profile.name}
                {!USER_SUBSCRIPTION_BADGE_MODE && Boolean(userBadge) && (
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
                  {profile.notes || USER_STATUS_MAP_TEXT[profile.status as number]}
                </Typography>
              )}
              <Typography variant="body2" component="p" color="textSecondary">
                <Typography component="span" variant="caption">
                  {profile.steam_id}
                </Typography>{' '}
                &middot; Joined {moment(profile.created_at).fromNow()}{' '}
                {moment().diff(moment(profile.created_at), 'days') <= USER_AGE_CAUTION && (
                  <NewbieIcon color="info" fontSize="inherit" sx={{ mb: -0.3 }} />
                )}
              </Typography>

              <Box sx={{ mb: 1 }}>
                <Typography variant="body2" component="span">
                  <Link href={`${linkProps.href}`}>{stats.live} Items</Link> &middot;{' '}
                  <Link href={`${linkProps.href}/reserved`}>{stats.reserved} Reserved</Link>{' '}
                  &middot; <Link href={`${linkProps.href}/delivered`}>{stats.sold} Delivered</Link>{' '}
                  &middot;{' '}
                  <Link href={`${linkProps.href}/bought`}>{stats.bid_completed} Bought</Link>
                </Typography>
                <Box sx={{ '& > *': { mt: 0.5 } }}>
                  {USER_SUBSCRIPTION_BADGE_MODE && userBadge && (
                    <>
                      <ExclusiveChip tag={userBadge} />
                      &nbsp;
                    </>
                  )}
                  {userTag && (
                    <>
                      <ExclusiveChip tag={userTag} />
                      &nbsp;
                    </>
                  )}
                  <ChipLink label="Steam Profile" href={profileURL} />
                  &nbsp;
                  <ChipLink label="Dotabuff" href={dotabuffURL} />
                </Box>
              </Box>
            </Box>
          </div>

          {isProfileReported ? (
            <Typography align="center">
              <br />
              <Button component={Link} nativeButton={false} href={`${linkProps.href}/activity`}>
                Show All Activity
              </Button>
            </Typography>
          ) : (
            <>
              <UserMarketList
                onSearchInput={handleSearchInput}
                data={markets}
                loading={loading}
                error={error}
                searchValue={String(qFilter ?? '')}
              />
              {!error && (
                <TablePaginationRouter
                  linkProps={linkProps}
                  style={{ textAlign: 'right' }}
                  count={markets.total_count}
                  page={filter.page as number}
                />
              )}
            </>
          )}
        </Container>

        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img src={trackProfileViewURL(profile.id)} height={1} width={1} alt="" />
      </main>

      <Footer />
    </>
  )
}

const marketSearchFilter = {
  page: 1,
  type: MARKET_TYPE_ASK,
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
  index: 'user_id',
}

// This gets called on every request
export async function getServerSideProps(context: {
  params: { id: string }
  query: Record<string, string | string[] | undefined>
}) {
  const { params, query } = context
  const vanityMode = Boolean(query.vanity)

  let profile: Profile | undefined
  let canonicalURL = ''

  // Check for vanity request.
  if (vanityMode) {
    try {
      profile = (await vanity(String(query.vanity))) as Profile
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
          profileError: (e as Error).message,
        },
      }
    }

    // When vanity exists use the profile from resolving it.
    // Otherwise try to get from users endpoint
  } else {
    try {
      profile = (await user(String(params.id))) as Profile
    } catch (e) {
      return {
        props: {
          profileError: (e as Error).message,
        },
      }
    }
  }

  // Retrieve initial user market summary.
  profile.stats = profile.market_stats

  // Retrieve initial user market data.
  let markets: MarketDatatable = {
    data: [],
    total_count: 0,
  }
  let marketError: string | null = null
  const filter: { page?: number; q?: string; user_id?: string | number } = {
    ...marketSearchFilter,
    user_id: profile.id,
  }
  filter.page = Number(query.page || 1)
  if (query.filter) {
    filter.q = String(query.filter)
  }

  try {
    const res = (await marketSearch(filter)) as MarketDatatable
    markets = res
  } catch (e) {
    marketError = (e as Error).message
  }

  // Compose profile page canonical URL.
  canonicalURL = `${APP_URL}/${vanityMode ? 'id' : 'profiles'}/${params.id}`

  // reduce large page data by nullifying un-needed data
  // nextjs.org/docs/messages/large-page-data
  // this can be reduce on API server level
  for (let i = 0; i < markets.data.length; i++) {
    const market = markets.data[i] as Partial<Market> & { inventory?: { steam_assets?: unknown } }
    delete market.user
    if (market.inventory) {
      delete market.inventory.steam_assets
    }
  }

  return {
    props: {
      profile,
      canonicalURL,
      filter,
      markets,
      profileError: null,
      marketError,
    },
  }
}
