import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Avatar from '@/components/Avatar'
import Typography from '@mui/material/Typography'
import { APP_NAME, APP_URL } from '@/constants/strings'
import { CDN_URL, marketSearch, user } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import MarketActivity from '@/components/MyMarketActivity'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
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

const defaultFilter = {
  sort: 'updated_at:desc',
  page: 1,
  limit: 15,
}

const defaultData = {
  data: [],
  total_result: 0,
  total_total: 0,
}

const scrollBias = 300

export default function UserActivity({ profile, canonicalURL }) {
  const { classes } = useStyles()

  const [datatable, setDatatable] = React.useState(defaultData)
  const [filter, setFilter] = React.useState({ ...defaultFilter, user_id: profile.id })
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState(null)

  React.useEffect(() => {
    if (loading) {
      return
    }

    setLoading(true)
    ;(async () => {
      try {
        const res = await marketSearch(filter)
        if (datatable.data.length === 0) {
          setDatatable(res)
        } else {
          const data = [...datatable.data, ...res.data]
          setDatatable({ ...datatable, data })
        }
      } catch (e) {
        setError(e.message)
      }
      setLoading(false)
    })()
  }, [filter])

  React.useEffect(() => {
    const listener = () => {
      const isLast = datatable.data.length === datatable.total_count
      if (loading || isLast || window.scrollY + scrollBias < window.scrollMaxY) {
        return
      }

      setFilter({ ...filter, page: filter.page + 1 })
    }

    window.addEventListener('scroll', listener)
    return () => {
      window.removeEventListener('scroll', listener)
    }
  })

  const profileURL = `/profiles/${profile.steam_id}`

  return (
    <>
      <Header />

      <Head>
        <meta charset="UTF-8" />
        <title>{`${APP_NAME} :: ${profile.name} items`}</title>
        <meta name="description" content={`${profile.name}'s delivered Giftable items`} />
        <link rel="canonical" href={canonicalURL} />
      </Head>

      <main className={classes.main}>
        <Container>
          <div>
            <Avatar
              className={classes.profile}
              src={`${CDN_URL}/${profile.avatar}`}
              glow={Boolean(profile.donation)}
              component={Link}
              href={profileURL}
            />
            <Typography variant="h6" color="textPrimary" component={Link} href={profileURL}>
              {profile.name}
            </Typography>
            <Typography color="textSecondary">
              {datatable && datatable.total_count} Market Activity
            </Typography>
          </div>
          <br />

          {error && <Typography color="error">{error.message.split(':')[0]}</Typography>}
          <MarketActivity datatable={datatable || {}} loading={loading} />
        </Container>
      </main>

      <Footer />
    </>
  )
}
UserActivity.propTypes = {
  profile: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
}

export async function getServerSideProps({ params }) {
  const profile = await user(String(params.id))
  const canonicalURL = `${APP_URL}/profiles/${params.id}/activity`

  return {
    props: {
      profile,
      canonicalURL,
    },
  }
}
