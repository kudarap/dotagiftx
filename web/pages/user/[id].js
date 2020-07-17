import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import moment from 'moment'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Link from '@material-ui/core/Link'
import Typography from '@material-ui/core/Typography'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { CDN_URL, marketSearch, user } from '@/service/api'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Container from '@/components/Container'
import UserMarketList from '@/components/UserMarketList'
import TablePagination from '@/components/TablePaginationRouter'

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

export default function UserDetails({ profile, markets, canonicalURL }) {
  const classes = useStyles()

  const router = useRouter()
  const [page, setPage] = React.useState(Number(router.query.page || 1))

  const handlePageChange = (e, p) => {
    setPage(p)
  }

  const linkProps = { href: '/user/[id]', as: `/user/${profile.steam_id}` }

  const profileURL = `https://steamcommunity.com/profiles/${profile.steam_id}`
  const steamRepURL = `https://steamrep.com/profiles/${profile.steam_id}`

  const metaTitle = `DotagiftX :: ${profile.name} profile`
  const metaDesc = `${profile.name}'s giftable items for sale`

  return (
    <>
      <Head>
        <title>Dota 2 Giftables :: {profile.name} listings</title>
        <meta name="description" content={`${profile.name} giftable listings`} />

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
                  href={steamRepURL}
                  variant="caption"
                  color="secondary"
                  target="_blank"
                  rel="noreferrer noopener">
                  {steamRepURL}
                </Link>
              </Typography>
            </Typography>
          </div>

          <UserMarketList data={markets} />
          <TablePagination
            linkProps={linkProps}
            style={{ textAlign: 'right' }}
            count={markets.total_count}
            page={page}
            onChangePage={handlePageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
UserDetails.propTypes = {
  profile: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
  markets: PropTypes.object,
}
UserDetails.defaultProps = {
  markets: {},
}

const marketSearchFilter = { status: MARKET_STATUS_LIVE, sort: 'created_at:desc' }

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query, req } = props
  const profile = await user(String(params.id))

  const filter = { ...marketSearchFilter, user_id: profile.id }
  if (query.page) {
    filter.page = Number(query.page)
  }

  const canonicalURL = req.headers.host + req.url

  return {
    props: {
      profile,
      canonicalURL,
      markets: await marketSearch(filter),
    },
  }
}
