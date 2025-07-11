import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Link from '@mui/material/Link'
import Avatar from '@/components/Avatar'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Button from '@/components/Button'
import SteamIcon from '@/components/SteamIcon'
import DiscordIcon from '@/components/DiscordIcon'
import { version } from '@/service/api'
import PropTypes from 'prop-types'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
    // background: 'url("/icon.png") no-repeat bottom right',
    // backgroundSize: 100,build
  },
}))

export default function About({ build }) {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: About</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            What about it?
          </Typography>
          <Typography component="h1" color="textSecondary">
            <Typography color="secondary" component="span" fontWeight="bold">
              {APP_NAME}
            </Typography>{' '}
            is short for Dota 2 giftables exchange, it was made to provide better search and pricing
            for Dota 2 giftable items like Collector&apos;s Caches which are not available on{' '}
            <Link href="https://steamcommunity.com" rel="noreferrer noopener" target="_blank">
              Steam Community Market
            </Link>
            . The project was heavily inspired by <strong>Giftable Megathread</strong> from{' '}
            <Link
              href="https://www.reddit.com/r/Dota2Trade"
              rel="noreferrer noopener"
              target="_blank">
              r/Dota2Trade
            </Link>
            .
          </Typography>
          <br />

          <Typography variant="h5" component="h2" gutterBottom>
            Who's behind it?
          </Typography>
          <Avatar src="/kudarap.jpg" style={{ width: 100, height: 100 }} />
          <Typography color="textSecondary">
            <strong>kudarap</strong> &mdash; author
            <br />
            Feel free to contact <strike>me</strike> us on Discord if you have issues or
            suggestions.
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
            disabled
            startIcon={<SteamIcon />}
            size="large"
            component={Link}
            target="_blank"
            rel="noreferrer noopener"
            href="https://steamcommunity.com/profiles/76561198088587178">
            Steam
          </Button>
          <Button
            startIcon={
              <img src="/icon_2x.png" style={{ height: 22, filter: 'brightness(10)' }} alt="dgx" />
            }
            size="large"
            component={Link}
            href="/profiles/76561198088587178">
            DotagiftX
          </Button>
          <br />
          <br />

          <Typography variant="h5" sx={{ mb: -1 }}>
            Version
          </Typography>
          <Typography color="textSecondary">
            <pre>
              tag: {build.version} <br />
              hash: {build.hash} <br />
              built: {build.built} <br />
            </pre>
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
About.propTypes = {
  build: PropTypes.object,
}
About.defaultProps = {
  build: {
    version: '-',
    hash: '-',
    built: '-',
  },
}

// This gets called on every request
export async function getServerSideProps() {
  // Fetch data from external API
  // const res = await fetch(API_URL)
  // const data = await res.json()
  const build = await version()

  // Pass data to the page via props
  return { props: { build } }
}
