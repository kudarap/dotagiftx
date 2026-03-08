import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Skeleton from '@mui/material/Skeleton'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import InternalUserCard from '@/components/InternalUserCard'
import { user } from '@/service/api'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  list: {
    listStyle: 'none',
    '& li:before': {
      content: `'🛡️ '`,
    },
    paddingLeft: theme.spacing(3),
    marginTop: 0,
  },
}))

const moderatorsUserIds = ['76561198078354099', '76561198171142718', '76561198057318750']

export default function Moderators() {
  const { classes } = useStyles()

  const [moderators, setModerators] = React.useState([])
  const [loading, setLoading] = React.useState(true)

  React.useEffect(() => {
    async function fetchUser(id) {
      return new Promise(resolve => {
        user(id).then(u => {
          resolve({
            id: u.steam_id,
            name: u.name,
            img: u.avatar,
            boons: u.boons,
            discordURL: 'https://discord.gg/UFt9Ny42kM',
          })
        })
      })
    }

    async function fetchUsers(array) {
      const res = await Promise.all(
        array.map(async item => {
          const v = await fetchUser(item)
          return v
        })
      )
      return res
    }

    fetchUsers(moderatorsUserIds).then(res => {
      setModerators(res)
      setLoading(false)
    })
  }, [])

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Moderators</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Moderators
            <Typography variant="body2" color="textSecondary">
              June 25, 2025
            </Typography>
          </Typography>
          <Typography color="textSecondary">
            The profiles listed below are the only official moderators of the site. Please head over
            to{' '}
            <Link href="https://discord.gg/UFt9Ny42kM" target="_blank" rel="noreferrer noopener">
              discord
            </Link>{' '}
            you need some questions or thoughts.
          </Typography>
          <br />

          {loading &&
            moderatorsUserIds.map(id => (
              <Skeleton key={id}>
                <InternalUserCard boons={['MODERATOR_TAG']} />
              </Skeleton>
            ))}

          {moderators.map(mod => (
            <InternalUserCard key={mod.id} {...mod} />
          ))}
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}

Moderators.propTypes = {
  users: PropTypes.arrayOf(PropTypes.object),
  error: PropTypes.string,
}
Moderators.defaultProps = {
  users: [],
  error: null,
}
