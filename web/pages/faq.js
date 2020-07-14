import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'

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
          <Typography variant="h5" component="h1" gutterBottom>
            Frequently Asked Questions
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            What is DotagiftX?
          </Typography>
          <Typography color="textSecondary">
            Market for Dota 2 giftables, items that can be gift or giftable-once are probably belong
            here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread.
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
            Why do you I need to sign with Steam?
          </Typography>
          <Typography color="textSecondary">
            It verifies Steam account ownership and provides some helpful links to check your
            profile and reputation.
          </Typography>
          <br />

          <Typography variant="h6" component="h2">
            Why do this?
          </Typography>
          <Typography color="textSecondary">
            Wanted to make tool that can be easily search and post these kind of items.
          </Typography>
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
