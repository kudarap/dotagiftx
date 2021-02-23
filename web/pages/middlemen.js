import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Link from '@material-ui/core/Link'
import Alert from '@material-ui/lab/Alert'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import ChipLink from '@/components/ChipLink'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

function Middleman({ name, id, internal = false }) {
  return (
    <>
      <strong>{name}</strong>
      {internal && (
        <>
          &nbsp;
          <ChipLink
            href={`https://dotagiftx.com/profiles/${id}`}
            label="DotagiftX"
            color="secondary"
          />
        </>
      )}
      &nbsp;
      <ChipLink href={`https://steamrep.com/profiles/${id}`} label="SteamRep" />
      &nbsp;
      <ChipLink href={`https://steamcommunity.com/profiles/${id}`} label="Steam Profile" />
    </>
  )
}

export default function Middlemen() {
  const classes = useStyles()

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Middlemen</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Alert severity="warning">
            This website is not responsible for scammed items and cannot help you recover them or
            help you scam the scammers.
          </Alert>
          <br />
          <br />

          {/* SteamRep middleman */}
          <Typography variant="h5" component="h1" gutterBottom>
            SteamRep&apos;s Official Middlemen
            <Typography variant="caption" color="textSecondary" component="sup">
              &nbsp;updated Feb 07 2021
            </Typography>
          </Typography>
          <Typography>
            <ul>
              <li>
                <Middleman name="kyuronite" id="76561198050680230" />
              </li>
              <li>
                <Middleman name="Hammy" id="76561197975564454" />
              </li>
              <li>
                <Middleman name="Eternal Mr Bones" id="76561198071974469" />
              </li>
              <li>
                <Middleman name="Alias" id="76561197982522773" />
              </li>
            </ul>

            <Typography color="textSecondary">
              Please <strong>READ</strong> and double check the&nbsp;
              <Link
                href="https://old.reddit.com/r/Dota2Trade/comments/l67zb1/psa_official_middlemen_have_green_shields_on/"
                target="_blank"
                color="secondary"
                rel="noreferrer noopener">
                source
              </Link>
              &nbsp;of the users listed above. You should see a green shield on their
              SteamRep&apos;s profile. In-case the list gets outdated you can check them on&nbsp;
              <Link
                href="https://reddit.com/r/Dota2Trade"
                target="_blank"
                color="secondary"
                rel="noreferrer noopener">
                r/Dota2Trade
              </Link>
              &nbsp;and&nbsp;
              <Link
                href="https://steamrep.com"
                target="_blank"
                color="secondary"
                rel="noreferrer noopener">
                SteamRep.com
              </Link>
            </Typography>
          </Typography>
          <br />
          <br />

          {/* DotagiftX middleman */}
          <Typography variant="h5" component="h2" gutterBottom>
            DotagiftX&apos;s Middleman
          </Typography>
          <ul>
            <li>
              <Middleman name="kudarap" id="76561198088587178" internal />
            </li>
          </ul>
          <Typography color="textSecondary">
            It&apos;s strongly recommended to get your middleman from Official SteamRep but if you
            trust
            <Link href="/profiles/76561198088587178" color="textPrimary">
              &nbsp;kudarap&nbsp;
            </Link>
            enough to middle your transaction, you can message a request on{' '}
            <Link
              href="https://discord.gg/UFt9Ny42kM"
              target="_blank"
              color="secondary"
              rel="noreferrer noopener">
              Discord
            </Link>{' '}
            to give a heads up.
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
