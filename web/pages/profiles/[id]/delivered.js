import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import useSWR from 'swr'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import { APP_NAME, APP_URL } from '@/constants/strings'
import { MARKET_STATUS_SOLD } from '@/constants/market'
import { CDN_URL, fetcher, MARKETS, statsMarketSummary, user } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
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

export default function UserDelivered({ profile, stats, canonicalURL }) {
  const classes = useStyles()

  filter.user_id = profile.id
  const { data, error, isValidating } = useSWR([MARKETS, filter], fetcher, {
    revalidateOnFocus: false,
  })

  const profileURL = `/profiles/${profile.steam_id}`

  return (
    <>
      <Header />

      <Head>
        <title>{`${APP_NAME} :: ${profile.name} delivered items`}</title>
        <meta name="description" content={`${profile.name}'s delivered Giftable items`} />
        <link rel="canonical" href={canonicalURL} />
      </Head>

      <main className={classes.main}>
        <Container>
          <div>
            <Avatar
              className={classes.profile}
              src={`${CDN_URL}/${profile.avatar}`}
              component={Link}
              href={profileURL}
            />
            <Typography variant="h6" color="textPrimary" component={Link} href={profileURL}>
              {profile.name}
            </Typography>
            <div style={{ display: 'flex' }}>
              <Typography component={Link} href={profileURL}>
                {stats.live} Items
              </Typography>
              &nbsp;&middot;&nbsp;
              <Typography component={Link} href={`${profileURL}/reserved`}>
                {stats.reserved} Reserved
              </Typography>
              &nbsp;&middot;&nbsp;
              <Typography
                component={Link}
                href={`${profileURL}/delivered`}
                style={{ textDecoration: 'underline' }}>
                {stats.sold} Delivered
              </Typography>
            </div>
          </div>
          <br />

          {error && <Typography color="error">{error.message.split(':')[0]}</Typography>}
          <MarketActivity datatable={data || {}} loading={isValidating} />
        </Container>
      </main>

      <Footer />
    </>
  )
}
UserDelivered.propTypes = {
  profile: PropTypes.object.isRequired,
  stats: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
}

export async function getServerSideProps({ params }) {
  const profile = await user(String(params.id))
  const canonicalURL = `${APP_URL}/profiles/${params.id}/delivered`

  const stats = await statsMarketSummary({ user_id: profile.id })

  return {
    props: {
      profile,
      canonicalURL,
      stats,
    },
  }
}
