import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Link from '@material-ui/core/Link'
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

function User({ name, id }) {
  return (
    <>
      <strong>{name}</strong>
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
          <Typography variant="h5" component="h1" gutterBottom>
            SteamRep&apos;s Official Middlemen
            <br />
            <Typography variant="caption" color="textSecondary">
              updated Feb 03 2021
            </Typography>
          </Typography>

          <Typography>
            <ul>
              <li>
                <User name="kyuronite" id="76561198050680230" />
              </li>
              <li>
                <User name="Hammy" id="76561197975564454" />
              </li>
              <li>
                <User name="Eternal Mr Bones" id="76561198071974469" />
              </li>
              <li>
                <User name="Alias" id="76561197982522773" />
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
              <br />
              <br />
              <strong>
                This website is not responsible for scammed items and cannot help you recover them.
              </strong>
            </Typography>
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
