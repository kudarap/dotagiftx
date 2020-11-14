import React from 'react'
import Link2 from 'next/link'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import Footer from '@/components/Footer'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(6),
  },
}))

export default function Faq() {
  const classes = useStyles()

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Link2 href="/profiles/76561198287849998" shallow>
            TEST SHALLOW
          </Link2>
          <Typography variant="h5" component="h1" gutterBottom>
            Frequently Asked Questions
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            What is DotagiftX?
          </Typography>
          <Typography color="textSecondary">
            Market for Dota 2 giftables, items that can be gift or giftable-once are probably belong
            here. If you are on Dota2Trade subreddit, its basicxy the Giftable Megathread.
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            What items can I find or post here?
          </Typography>
          <Typography color="textSecondary">
            Anything Dota 2 items that can be gift to a friend like items from Collector&apos;s
            Caches, In-game drops, or Immortal treasures.
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            Why do you I need to sign in with Steam?
          </Typography>
          <Typography color="textSecondary">
            It verifies Steam account ownership and provides some helpful links to check your
            profile and reputation.
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            Can I trust the users on this website?
          </Typography>
          <Typography color="textSecondary">
            No, but there are quick links for SteamRep and Steam Profile to help you validate them.
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            How can I report scammers?
          </Typography>
          <Typography color="textSecondary">
            You can use <Link href="https://steamrep.com/">SteamRep</Link> or inquire on{' '}
            <Link href="https://www.reddit.com/r/Dota2Trade/">r/Dota2Trade</Link>.
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            Why do this?
          </Typography>
          <Typography color="textSecondary">
            Wanted to make tool that can be easily search and post these giftable items.
          </Typography>
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
