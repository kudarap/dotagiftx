import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import useSWR from 'swr'
import moment from 'moment'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import Divider from '@material-ui/core/Divider'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { Link } from '@material-ui/core'
import { CDN_URL, fetcher, MARKETS, user } from '@/service/api'
import { APP_URL } from '@/constants/strings'
import {
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_RESERVED,
} from '@/constants/market'
import { dateFromNow } from '@/lib/format'
import ItemImage from '@/components/ItemImage'

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
            <Avatar src={`${CDN_URL}/${profile.avatar}`} />
            Reserved Items
          </Typography>
          {error && <Typography color="error">{error.message.split(':')[0]}</Typography>}
          <ul style={{ paddingLeft: 0, listStyle: 'none' }}>
            {!isValidating &&
              data &&
              data.data.map(market => (
                <li>
                  <ItemImage
                    style={{ width: 60, height: 40, marginRight: 8, float: 'left' }}
                    image={`/200x100/${market.item.image}`}
                    title={market.item.name}
                    rarity={market.item.rarity}
                  />
                  <Typography variant="body2">
                    {market.user.name}{' '}
                    <span style={{ color: MARKET_STATUS_MAP_COLOR[market.status] }}>
                      {MARKET_STATUS_MAP_TEXT[market.status].toLowerCase()}
                    </span>
                    &nbsp;
                    <Link href={`/item/${market.item.slug}`} color="secondary">
                      {market.item.name}
                    </Link>
                    &nbsp;
                    {moment(market.updated_at).fromNow()}
                  </Typography>
                  <Typography component="pre" color="textSecondary" variant="caption">
                    {market.notes}
                  </Typography>
                  <Divider style={{ margin: '8px 0 8px' }} light />
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
