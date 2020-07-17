import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import moment from 'moment'
import useSWR from 'swr'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Link from '@material-ui/core/Link'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { CDN_URL, MARKETS, marketSearch, user, fetcher } from '@/service/api'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Container from '@/components/Container'
import UserMarketList from '@/components/UserMarketList'

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

const marketSearchFilter = {
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
}

export default function UserDetails({ profile, markets: initialData }) {
  const classes = useStyles()

  marketSearchFilter.user_id = profile.id
  const { data: marketListing, error: marketError } = useSWR(
    [MARKETS, marketSearchFilter],
    (u, f) => fetcher(u, f),
    { initialData }
  )

  const profileURL = `https://steamcommunity.com/profiles/${profile.steam_id}`
  const steamrepURL = `https://steamrep.com/profiles/${profile.steam_id}`

  return (
    <>
      <Head>
        <title>Dota 2 Giftables :: {profile.name} listings</title>
        <meta name="description" content={`${profile.name} giftable listings`} />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <div className={classes.details}>
            <Avatar className={classes.avatar} src={CDN_URL + profile.avatar} />
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
                  {`steam: `}
                </Typography>
                <Link
                  href={profileURL}
                  variant="caption"
                  color="secondary"
                  target="_blank"
                  rel="noreferrer noopener">
                  {profileURL}
                </Link>
                <br />

                <Typography color="textSecondary" component="span">
                  {`steamrep: `}
                </Typography>
                <Link
                  href={steamrepURL}
                  variant="caption"
                  color="secondary"
                  target="_blank"
                  rel="noreferrer noopener">
                  {steamrepURL}
                </Link>
              </Typography>
            </Typography>
          </div>

          <UserMarketList data={marketListing} error={marketError} />
        </Container>
      </main>

      <Footer />
    </>
  )
}
UserDetails.propTypes = {
  profile: PropTypes.object.isRequired,
  markets: PropTypes.object,
}
UserDetails.defaultProps = {
  markets: {},
}

// This gets called on every request
export async function getServerSideProps({ params }) {
  const profile = await user(String(params.id))

  marketSearchFilter.user_id = profile.id
  return {
    props: {
      profile,
      markets: await marketSearch(marketSearchFilter),
    },
  }
}
