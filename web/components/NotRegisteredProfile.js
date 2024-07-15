import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Avatar from '@/components/Avatar'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Container from '@/components/Container'
import ChipLink from '@/components/ChipLink'
import {
  APP_NAME,
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import { Alert } from '@mui/material'
import Link from '@/components/Link'
import { dateFromNow } from '@/lib/format'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  profileName: {
    [theme.breakpoints.down('sm')]: {
      fontSize: theme.typography.h6.fontSize,
    },
  },
  details: {
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  avatar: {
    [theme.breakpoints.down('sm')]: {
      margin: '0 auto',
    },
    width: 100,
    height: 100,
    marginRight: theme.spacing(1.5),
  },
}))

export default function NotRegisteredProfile({ profile, canonicalURL }) {
  const { classes } = useStyles()

  const profileURL = `${STEAM_PROFILE_BASE_URL}/${profile.steam_id}`
  const steamRepURL = `${STEAMREP_PROFILE_BASE_URL}/${profile.steam_id}`
  const dotabuffURL = `${DOTABUFF_PROFILE_BASE_URL}/${profile.steam_id}`

  const metaTitle = `${APP_NAME} :: ${profile.name}`
  const metaDesc = `${profile.name}'s Steam profile`

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={canonicalURL} />

        {/* Twitter Card */}
        <meta name="twitter:card" content="summary" />
        <meta name="twitter:title" content={metaTitle} />
        <meta name="twitter:description" content={metaDesc} />
        <meta name="twitter:image" content={profile.steam_avatar} />
        <meta name="twitter:site" content={`@${APP_NAME}`} />
        {/* OpenGraph */}
        <meta property="og:url" content={canonicalURL} />
        <meta property="og:type" content="website" />
        <meta property="og:title" content={metaTitle} />
        <meta property="og:description" content={metaDesc} />
        <meta property="og:image" content={profile.steam_avatar} />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <div className={classes.details}>
            <Avatar className={classes.avatar} src={profile.steam_avatar} />
            <Typography component="h1">
              <Typography className={classes.profileName} component="p" variant="h4">
                {profile.name}
              </Typography>
              <Typography gutterBottom>
                <Link href={profile.url} variant="body2">
                  {profile.url}
                </Link>
                <br />
                <ChipLink label="Steam Profile" href={profileURL} />
                &nbsp;
                {/* <ChipLink label="Steam Inventory" href={`${profileURL}/inventory`} /> */}
                {/* &nbsp; */}
                <ChipLink label="SteamRep" href={steamRepURL} />
                &nbsp;
                <ChipLink label="Dotabuff" href={dotabuffURL} />
              </Typography>
            </Typography>
          </div>
          <Typography gutterBottom />

          <Alert severity="info" variant="outlined">
            This Steam user is not registered on DotagiftX. Profile data updated{' '}
            {dateFromNow(profile.last_updated_at)}
          </Alert>
        </Container>
      </main>

      <Footer />
    </>
  )
}
NotRegisteredProfile.propTypes = {
  profile: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
}
