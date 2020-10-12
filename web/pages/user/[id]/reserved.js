import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import useSWR from 'swr'
import moment from 'moment'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { Link } from '@material-ui/core'
import { CDN_URL, fetcher, MARKETS, user } from '@/service/api'
import { APP_URL } from '@/constants/strings'
import { MARKET_STATUS_MAP_TEXT, MARKET_STATUS_RESERVED } from '@/constants/market'
import { dateFromNow } from '@/lib/format'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
}))

const filter = {
  status: MARKET_STATUS_RESERVED,
  sort: 'updated_at:desc',
}

export default function UserReserved({ profile, canonicalURL }) {
  const classes = useStyles()

  filter.user_id = profile.id
  const { data, error, isValidating } = useSWR([MARKETS, filter], fetcher)

  return (
    <>
      <Header />

      <Head>
        <title>{`DotagiftX :: ${profile.name} reserved items`}</title>
        <meta name="description" content={`${profile.name}'s on-reserved giftable items`} />
        <link rel="canonical" href={canonicalURL} />
      </Head>

      <main className={classes.main}>
        <Container>
          <Typography component="h1" gutterBottom>
            Reserved Items
          </Typography>
          <ul>
            {data &&
              data.data.map(market => (
                <li>
                  <Typography variant="body2">
                    {market.user.name} {MARKET_STATUS_MAP_TEXT[market.status].toLowerCase()}&nbsp;
                    <Link href={`/item/${market.item.slug}`} color="secondary">
                      {market.item.name}
                    </Link>
                    &nbsp;
                    {moment(market.updated_at).fromNow()}
                  </Typography>
                  <Typography
                    component="pre"
                    color="textSecondary"
                    variant="caption"
                    paragraph
                    style={{ marginLeft: 4 }}>
                    {market.notes}
                  </Typography>
                </li>
              ))}
          </ul>
        </Container>
      </main>

      <Footer />
    </>
  )
}
UserReserved.propTypes = {
  profile: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
}

export async function getServerSideProps({ params }) {
  const profile = await user(String(params.id))
  const canonicalURL = `${APP_URL}/user/${params.id}/reserve`

  return {
    props: {
      profile,
      canonicalURL,
    },
  }
}
