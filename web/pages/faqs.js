import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import MuiLink from '@mui/material/Link'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import Footer from '@/components/Footer'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  question: {
    paddingTop: theme.spacing(2.5),
    '&:target': {
      borderBottom: `2px inset ${theme.palette.secondary.main}`,
      '& .MuiLink-root:hover': {
        textDecoration: 'none',
      },
    },
  },
}))

function slugify(s) {
  return String(s)
    .toLowerCase()
    .replace(/[^a-z0-9 -]/g, '')
    .replace(/\s+/g, '-')
}

function Question({ children, ...other }) {
  const { classes } = useStyles()
  const id = slugify(children)
  return (
    <Typography
      className={classes.question}
      component="h2"
      id={id}
      gutterBottom
      style={{ fontWeight: 'bold' }}
      {...other}>
      <MuiLink href={`#${id}`} color="textPrimary">
        {children}
      </MuiLink>
    </Typography>
  )
}
function Answer({ children }) {
  return (
    <Typography color="textSecondary" gutterBottom>
      {children}
    </Typography>
  )
}

export default function Faqs() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Frequently Asked Questions</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Frequently Asked Questions
          </Typography>

          <Question>What is DotagiftX?</Question>
          <Answer>
            Market place for Dota 2 Giftables, items that can only be gift or gift-once are probably
            belong here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread
            with a kick.
          </Answer>

          <Question>What items I can find or post here?</Question>
          <Answer>
            Anything Dota 2 items that can be gift to a friend like set bundles from
            Collector&apos;s Cache, Immortal Treasure items, or rare in-game drops from Treasure of
            the Cryptic Beacon.
          </Answer>

          <Question>Why do I need to sign in with Steam?</Question>
          <Answer>
            You don&apos;t, if you are looking for offers it&apos;s open to public and you can
            contact the seller. If you are listing your items or placing a request we need to verify
            Steam account ownership and it will help users to check your profile and reputation.
          </Answer>

          <Question>Can I trust the users on this website?</Question>
          <Answer>
            Not really, its open for anyone so please be vigilant to scammers. User&apos;s
            transaction history are open and links to their SteamRep, Steam, and Dotabuff are listed
            for you to checkout.
          </Answer>

          <Question>What is reservation / reservation fee / deposit?</Question>
          <Answer>
            Some sellers requires a small fee to lock the item to buyer and this varies depending on
            the seller or the item. <em>Reserved item</em> status will not appear on search and
            profile listings to stop offering to other buyers. If you signed up, you can check your{' '}
            <Link color="secondary" href="/my-orders#toreceive">
              reservations here
            </Link>
            .
          </Answer>

          <Question>Why do I need to wait 30 days to send or receive an item?</Question>
          <Answer>
            Valve&apos;s gift restriction that you need to have 30 days as friend to send and
            receive Giftable items.
          </Answer>

          <Question>Do I need a Middleman?</Question>
          <Answer>
            If you asked, you probably do, specially on high-value items where scammers fuck around.
            DotagiftX ONLY suggest that you get the{' '}
            <Link href="/middlemen" color="secondary">
              Middleman here
            </Link>
            &nbsp; and read around.
          </Answer>

          <Question>How do I report scammers?</Question>
          <Answer>
            Please{' '}
            <Link
              href="https://discord.gg/UFt9Ny42kM"
              target="_blank"
              rel="noreferrer noopener"
              color="secondary">
              contact kudarap
            </Link>{' '}
            to ban the account on this site and you can submit a report on{' '}
            <Link
              href="https://steamrep.com/"
              target="_blank"
              color="secondary"
              rel="noreferrer noopener">
              SteamRep
            </Link>{' '}
            or inquire on{' '}
            <Link
              href="https://www.reddit.com/r/Dota2Trade/"
              target="_blank"
              color="secondary"
              rel="noreferrer noopener">
              r/Dota2Trade
            </Link>
            .
          </Answer>

          <Question>Why do this?</Question>
          <Answer>
            Wanted to sell Giftable items using a website so it can be googled, and might be useful
            to others.
            {/* Wanted to make tool that can be easily search and post these Giftable items. */}
          </Answer>
        </Container>
      </main>

      <Footer />
    </>
  )
}
