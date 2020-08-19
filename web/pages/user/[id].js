import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import moment from 'moment'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { CDN_URL, marketSearch, user } from '@/service/api'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Container from '@/components/Container'
import ChipLink from '@/components/ChipLink'
import UserMarketList from '@/components/UserMarketList'
import TablePaginationRouter from '@/components/TablePaginationRouter'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
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

const steamProfileBaseURL = 'https://steamcommunity.com/profiles'
const steamRepBaseURL = 'https://steamrep.com/profiles'

export default function UserDetails({
  profile,
  filter,
  markets: initialMarkets,
  error: initialError,
  canonicalURL,
}) {
  const classes = useStyles()

  const [page, setPage] = React.useState(filter.page)
  const [markets, setMarkets] = React.useState(initialMarkets)
  const [error, setError] = React.useState(initialError)

  // Handle market request on page change.
  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await marketSearch({ ...filter, page })
        setMarkets(res)
      } catch (e) {
        setError(e.message)
      }
    })()
  }, [page])

  const handlePageChange = (e, p) => {
    setPage(p)
  }

  const linkProps = { href: '/user/[id]', as: `/user/${profile.steam_id}` }

  const profileURL = `${steamProfileBaseURL}/${profile.steam_id}`
  const steamRepURL = `${steamRepBaseURL}/profiles/${profile.steam_id}`

  const metaTitle = `DotagiftX :: ${profile.name}`
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
        <meta name="twitter:site" content="@DotagiftX" />
        {/* OpenGraph */}
        <meta property="og:url" content={canonicalURL} />
        <meta property="og:type" content="article" />
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
              <Typography component="p" variant="h4">
                {profile.name}
              </Typography>
              <Typography gutterBottom>
                <Typography color="textSecondary" component="span">
                  {`registered: `}
                </Typography>
                {moment(profile.created_at).fromNow()}
                <br />
                <Typography color="textSecondary" component="span">
                  {`quick links: `}
                </Typography>
                <ChipLink label="Steam Profile" href={profileURL} />
                &nbsp;
                {/* <ChipLink label="Steam Inventory" href={`${profileURL}/inventory`} /> */}
                {/* &nbsp; */}
                <ChipLink label="SteamRep" href={steamRepURL} />
              </Typography>
            </Typography>
          </div>

          <UserMarketList data={markets} error={error} />
          {!error && (
            <TablePaginationRouter
              linkProps={linkProps}
              style={{ textAlign: 'right' }}
              count={markets.total_count}
              page={page}
              onChangePage={handlePageChange}
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
export async function getServerSideProps(props) {
  const { params, query, req } = props

  const canonicalURL = `https://${req.headers.host}${req.url}`

  const profile = await user(String(params.id))
  const filter = { ...marketSearchFilter, user_id: profile.id }
  if (query.page) {
    filter.page = Number(query.page)
  }

  let markets = {}
  let error = null
  try {
    markets = await marketSearch(filter)
  } catch (e) {
    error = e.message
  }

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
