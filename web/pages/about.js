import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@/components/Avatar'
import Typography from '@material-ui/core/Typography'
import Link from '@material-ui/core/Link'
import RedditIcon from '@material-ui/icons/Reddit'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Button from '@/components/Button'
// import SteamIcon from '@/components/SteamIcon'
import DiscordIcon from '@/components/DiscordIcon'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
    // background: 'url("/icon.png") no-repeat bottom right',
    // backgroundSize: 100,
  },
}))

export default function About() {
  const classes = useStyles()

  return (
    <>
      <Head>
        <title>{APP_NAME} :: About</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Who is behind this?
          </Typography>
          <br />
          <Avatar src="/kudarap.jpg" style={{ width: 100, height: 100 }} />
          <Typography color="textSecondary">
            <strong>kudarap</strong> &mdash; author
            <br />
            Feel free to contact me if you have issues or suggestions.
          </Typography>
          <Button
            startIcon={<DiscordIcon />}
            size="large"
            component={Link}
            target="_blank"
            rel="noreferrer noopener"
            href="https://discord.gg/UFt9Ny42kM">
            Discord
          </Button>
          <Button
            startIcon={<RedditIcon />}
            size="large"
            component={Link}
            target="_blank"
            rel="noreferrer noopener"
            href="https://www.reddit.com/message/compose/?to=kudarap">
            Reddit
          </Button>
          {/* <Button
            startIcon={<SteamIcon />}
            size="large"
            component={Link}
            target="_blank"
            rel="noreferrer noopener"
            href="https://steamcommunity.com/profiles/76561198088587178">
            Steam
          </Button> */}
          <Button
            startIcon={
              <img src="/icon_2x.png" style={{ height: 22, filter: 'brightness(10)' }} alt="dgx" />
            }
            size="large"
            component={Link}
            href="/profiles/76561198088587178">
            DotagiftX
          </Button>
        </Container>
      </main>

      <Footer />
    </>
  )
}
