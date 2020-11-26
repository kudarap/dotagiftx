import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import Footer from '@/components/Footer'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(6),
  },
}))

function Question({ children }) {
  return (
    <Typography component="h2">
      <strong>{children}</strong>
    </Typography>
  )
}

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

          <Question>What is {APP_NAME}?</Question>
          <Typography color="textSecondary">
            Market place for Dota 2 Giftables, items that can only be gift or gift-once are probably
            belong here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread
            with a kick.
          </Typography>
          <br />

          <Question>What items I can find or post here?</Question>
          <Typography color="textSecondary">
            Anything Dota 2 items that can be gift to a friend like set bundles from
            Collector&apos;s Cache, In-game drops, or Immortal treasures.
          </Typography>
          <br />

          <Question>Why do I need to sign in with Steam?</Question>
          <Typography color="textSecondary">
            It verifies Steam account ownership and provides some helpful links to check your
            profile and reputation.
          </Typography>
          <br />

          <Question>Can I trust the users on this website?</Question>
          <Typography color="textSecondary">
            Not really, but there are quick links like SteamRep and Steam on their profile to help
            you check them.
          </Typography>
          <br />

          <Question>Why do I need to wait 30 days to send the item?</Question>
          <Typography color="textSecondary">
            Valve&apos;s rule that you need to have 30-day cooldown as friend to send giftable
            items.
          </Typography>
          <br />

          <Question>How do I report scammers?</Question>
          <Typography color="textSecondary">
            You can use <Link href="https://steamrep.com/">SteamRep</Link> or inquire on{' '}
            <Link href="https://www.reddit.com/r/Dota2Trade/">r/Dota2Trade</Link>.
          </Typography>
          <br />

          <Question>Why do this?</Question>
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
