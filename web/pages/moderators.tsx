import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Skeleton from '@mui/material/Skeleton'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import InternalUserCard from '@/components/InternalUserCard'
import type { InternalUserCardProps } from '@/components/InternalUserCard'
import { user } from '@/service/api'
import Link from '@/components/Link'
import type { Profile } from '@/lib/types'

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

  const [moderators, setModerators] = React.useState<InternalUserCardProps[]>([])
  const [loading, setLoading] = React.useState(true)

  React.useEffect(() => {
    async function fetchUser(id: string) {
      return new Promise(resolve => {
        user(id).then(u => {
          const profile = u as Profile
          resolve({
            id: profile.steam_id,
            name: profile.name,
            img: profile.avatar,
            boons: profile.boons,
            discordURL: 'https://discord.gg/UFt9Ny42kM',
            createdAt: profile.created_at,
          })
        })
      })
    }

    async function fetchUsers(array: string[]) {
      const res = await Promise.all(
        array.map(async item => {
          const v = await fetchUser(item)
          return v
        })
      )
      return res
    }

    fetchUsers(moderatorsUserIds).then(res => {
      setModerators(res as InternalUserCardProps[])
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
