import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import useSWR from 'swr'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { Link } from '@material-ui/core'
import { CDN_URL, fetcher, MARKETS, user } from '@/service/api'
import { APP_URL } from '@/constants/strings'
import { MARKET_STATUS_SOLD } from '@/constants/market'
import MarketActivity from '@/components/MarketActivity'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
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

const filter = {
  status: MARKET_STATUS_SOLD,
  sort: 'updated_at:desc',
  limit: 50,
}

export default function UserDelivered({ profile, canonicalURL }) {
  const classes = useStyles()

  filter.user_id = profile.id
  const { data, error, isValidating } = useSWR([MARKETS, filter], fetcher)

  return (
    <>
      <Header />

      <Head>
        <title>{`DotagiftX :: ${profile.name} delivered items`}</title>
        <meta name="description" content={`${profile.name}'s delivered giftable items`} />
        <link rel="canonical" href={canonicalURL} />
      </Head>

      <main className={classes.main}>
        <Container>
          <div>
            <Avatar
              className={classes.profile}
              src={`${CDN_URL}/${profile.avatar}`}
              component={Link}
              href={`/profiles/${profile.steam_id}`}
            />
            <Typography
              variant="h6"
              color="textPrimary"
              component={Link}
              href={`/profiles/${profile.steam_id}`}>
              {profile.name}
            </Typography>
            <Typography color="textSecondary">
              {data && data.total_count} Delivered Items
            </Typography>
          </div>
          {error && <Typography color="error">{error.message.split(':')[0]}</Typography>}
          <MarketActivity data={data ? data.data : null} loading={isValidating} />
        </Container>
      </main>

      <Footer />
    </>
  )
}
UserDelivered.propTypes = {
  profile: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
}

export async function getServerSideProps({ params }) {
  const profile = await user(String(params.id))
  const canonicalURL = `${APP_URL}/profiles/${params.id}/reserve`

  return {
    props: {
      profile,
      canonicalURL,
    },
  }
}
