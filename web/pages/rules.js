import React from 'react'
import Head from 'next/head'
import Typography from '@mui/material/Typography'
import { makeStyles } from 'tss-react/mui'

import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Link from '@/components/Link'
import { APP_NAME } from '@/constants/strings'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  ruleContainer: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
    padding: theme.spacing(2),
    background: '#190f00',
    color: '#c1b79b',
    borderRadius: 5,
  },
}))

export default function Version() {
  const { classes } = useStyles()

  return (
    <div className="container">
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Rules</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom style={{ marginBottom: 26 }}>
            Rules
            <Typography variant="body2" color="textSecondary">
              Sep 19, 2021
            </Typography>
          </Typography>
          <Typography color="textSecondary">
            To keep the community fair and a little bit safe for both sellers and buyers. Here are
            some written rules to follow.
          </Typography>

          <div className={classes.ruleContainer}>
            <Typography component="h2" style={{ fontWeight: 'bold' }}>
              Selling Rules
            </Typography>

            <Typography>
              <ol>
                <li>No misleading post or abusive behavior.</li>
                <li>No multiple reservations on a single item or stock.</li>
                <li>No reservation cancellation without prior notice.</li>
              </ol>
            </Typography>

            <Typography>
              Breaking one of the rules will be punishable by <u>Account Suspension</u> to{' '}
              <u>Permanent Ban</u>.
            </Typography>
          </div>

          <div className={classes.ruleContainer}>
            <Typography component="h2" style={{ fontWeight: 'bold' }}>
              Buying Rules
            </Typography>

            <Typography>
              <ol>
                <li>No reservation cancellation without prior notice.</li>
              </ol>
            </Typography>

            <Typography>
              Breaking the rule will reflect a negative impact on buyer&apos;s profile.
            </Typography>
          </div>

          <Typography color="textSecondary">
            Since we don&apos;t have profile feedback system yet. You can report them to{' '}
            <Link href="/feedback" color="secondary">
              Feedback
            </Link>{' '}
            page along with details. Please don&apos;t forget to include their dotagiftx profile and
            if you have images please use link from imgur and alike.
            <br />
            <br />
            If you have comments or suggestions regarding this feel free to reach out on{' '}
            <Link
              href="https://discord.gg/UFt9Ny42kM"
              target="_blank"
              color="secondary"
              rel="noreferrer noopener">
              Discord
            </Link>{' '}
            or submit a comment on feedback.
          </Typography>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
