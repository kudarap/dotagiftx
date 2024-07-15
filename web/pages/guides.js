import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function Blacklist() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Guides</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Guides & Tips
          </Typography>
          <br />

          <Typography style={{ fontWeight: 'bold' }} gutterBottom>
            Profile checking
            <Typography color="textSecondary">
              You can go to their Steam profiles and change the url{' '}
              <Typography
                component="em"
                color="textPrimary"
                variant="body2"
                style={{ fontWeight: 500 }}>
                steamcommunity.com
              </Typography>{' '}
              to{' '}
              <Typography
                component="em"
                color="textPrimary"
                variant="body2"
                style={{ fontWeight: 500 }}>
                dotagiftx.com
              </Typography>{' '}
              that will lead you to their Dotagiftx profile with transaction history and links to
              SteamRep and Dotabuff. You can search scammers on{' '}
              <Link href="/bans" color="secondary">
                Banned users
              </Link>
              .
            </Typography>
          </Typography>
          <br />

          <Typography style={{ fontWeight: 'bold' }} gutterBottom>
            Before going first
            <Typography color="textSecondary">
              After waiting for 30days and before going first, its highly recommend that you check
              the profile again for scam alerts or reports. Most of the victim scam does not know
              their seller was already reported and could have been avoided.
            </Typography>
          </Typography>
          <br />

          <Typography style={{ fontWeight: 'bold' }} gutterBottom>
            Prepare evidences
            <Typography color="textSecondary">
              Prepare in-case the seller scam you. Take screenshots of things during transaction
              that you can use to submit a case on SteamRep. If you buying here on DotagiftX make
              sure that you have a record of reservation so we could track the transaction.
            </Typography>
          </Typography>
          <br />

          <Typography style={{ fontWeight: 'bold' }}>
            Gift restrictions
            <Typography color="textSecondary">
              To combat credit card fraud, Valve has placed several restrictions on gifting.{' '}
              <Link
                href="https://dota2.fandom.com/wiki/Wrapped_Gift"
                target="_blank"
                color="secondary"
                rel="noreferrer noopener">
                dota2.fandom.com
              </Link>
            </Typography>
          </Typography>
          <Typography color="textSecondary">
            <ul>
              <li>Gifts can be sent to a friend, but not sold on the Steam Market</li>
              <li>
                Gifts can only be sent to friends of at least one year, unless Mobile Authenticator
                is activated.
              </li>
              <li>
                Only a limited number of items can be gifted within a certain time. ( 8 items for
                every 24 hours )
                <li>
                  The number of gifts one can send is determined by the player's Experience Trophy.
                </li>
              </li>
            </ul>
          </Typography>
          <br />

          <Typography style={{ fontWeight: 'bold' }}>Buying giftables</Typography>
          <Typography color="textSecondary">
            <ul>
              <li>Always check the item or set availability on seller&apos;s Dota 2 inventory.</li>
              <li>
                Dota 2 Giftables transaction only viable if the two steam user parties have been
                friends for 30 days.
              </li>
              <li>
                As Giftables involves a party having to go first, please always check seller&apos;s
                reputation through&nbsp; SteamRep and transaction history.
              </li>

              <li>
                If you need a middleman, I only suggest you get{' '}
                <Link href="/middlemen" color="secondary">
                  Middleman here
                </Link>
                .
              </li>
            </ul>
          </Typography>
          <br />

          <Typography style={{ fontWeight: 'bold' }}>Selling giftables</Typography>
          <Typography color="textSecondary">
            <ul>
              <li>Please be respectful on the price stated by the buyer.</li>
              <li>Make sure your item exist in your inventory.</li>
              <li>
                Dota 2 Giftables transaction only viable if the two steam user parties have been
                friends for 30 days.
              </li>
              <li>
                Payment agreements will be done between you and the buyer. This website does not
                accept or integrate any payment service.
              </li>
            </ul>
          </Typography>

          <Typography>
            You can read more on{' '}
            <Link href="/faqs" color="secondary">
              FAQs
            </Link>
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
