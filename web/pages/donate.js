import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import AwardIcon from '@mui/icons-material/Flare'
import KeyIcon from '@mui/icons-material/VpnKey'
import MoneyIcon from '@mui/icons-material/LocalAtm'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Link from '@/components/Link'
import Button from '@/components/Button'
import ProfileCard from '@/components/ProfileCard'
import Table from '@mui/material/Table'
import { Paper } from '@mui/material'
import DonatorBadge from '@/components/DonatorBadge'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  item: {
    display: 'flex',
    alignItems: 'center',
  },
  icon: {
    marginRight: theme.spacing(2),
  },
}))

export default function Faq() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Donate</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            How can I donate?
            <Typography color="textSecondary">
              Thanks for your interest supporting DotagiftX! Since its not monetizing on views and I
              do not run ads on this site, giving a{' '}
              <Link href="https://steamcommunity.com/sharedfiles/filedetails/?id=2313234224">
                thumbs up on Steam
              </Link>{' '}
              or giving feedback is good enough to help the site for the time being. You can check
              below how you can help in other ways.
              {/* Currently this site is running on someone else&apos;s server that is why you don&apos;t see ads */}
            </Typography>
          </Typography>
          <br />

          <div className={classes.item}>
            <AwardIcon className={classes.icon} color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Give the project a{' '}
              <Typography color="textPrimary" component="span">
                Steam Award
              </Typography>{' '}
              or thumbs up on Steam is very much appreciated.
              <Button
                target="_blank"
                rel="noreferrer noopener"
                href="https://steamcommunity.com/sharedfiles/filedetails/?id=2313234224"
                color="secondary"
                size="small">
                Give Award
              </Button>
            </Typography>
          </div>
          <br />

          <div className={classes.item}>
            <KeyIcon className={classes.icon} color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Donate tradable keys (
              <Typography color="textPrimary" component="span">
                TF2
              </Typography>{' '}
              or{' '}
              <Typography color="textPrimary" component="span">
                CS:GO case keys
              </Typography>
              ) and will give you a{' '}
              <span style={{ textDecoration: 'line-through' }}>donator badge</span> on your profile
              in return.
              <Button
                target="_blank"
                rel="noreferrer noopener"
                href="https://steamcommunity.com/tradeoffer/new/?partner=128321450&token=38BJlyuW"
                color="secondary"
                size="small">
                Send Trade offer
              </Button>
            </Typography>
          </div>
          <br />

          <div className={classes.item}>
            <MoneyIcon className={classes.icon} color="inherit" fontSize="large" />
            <Typography color="textSecondary">
              Toss a coin on{' '}
              <Typography color="textPrimary" component="span">
                Paypal
              </Typography>{' '}
              or{' '}
              <Typography color="textPrimary" component="span">
                Crypto
              </Typography>{' '}
              if you want to help future cost of server and domain. Please don&apos;t forget to put
              your profile link on the notes so I can award the{' '}
              <span style={{ textDecoration: 'line-through' }}>donator badge</span>.
              <Button
                component={Link}
                href="/donate-coin"
                color="secondary"
                size="small"
                underline="none">
                Donate Coin
              </Button>
            </Typography>
          </div>
          <br />

          <Typography
            color="textSecondary"
            component="em"
            style={{ textDecoration: 'line-through' }}>
            Donator badge will be on your profile forever and make your avatar glow for 30 days.
          </Typography>

          {/*<br />*/}

          {/*<Typography variant="h5">What the badge looks like?</Typography>*/}
          {/*<Paper style={{ padding: 24, margin: '10px auto', width: 600 }}>*/}
          {/*  <ProfileCard*/}
          {/*    user={{*/}
          {/*      id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',*/}
          {/*      steam_id: '76561198088587178',*/}
          {/*      name: 'kudarap',*/}
          {/*      url: 'https://steamcommunity.com/id/kudarap/',*/}
          {/*      avatar: '6401d0c455e255525e605d328b66375099e46bb2.jpg',*/}
          {/*      status: 0,*/}
          {/*      donation: 1,*/}
          {/*      donated_at: '2021-04-22T04:22:16.613Z',*/}
          {/*      created_at: '2020-06-18T13:13:43.926+08:00',*/}
          {/*      updated_at: '2021-03-16T14:23:59.871+08:00',*/}
          {/*    }}*/}
          {/*    marketSummary={{*/}
          {/*      live: 92,*/}
          {/*      reserved: 2,*/}
          {/*      sold: 26,*/}
          {/*    }}*/}
          {/*  />*/}
          {/*</Paper>*/}
        </Container>
      </main>

      <Footer />
    </>
  )
}
