import React from 'react'
import Head from 'next/head'
import Typography from '@material-ui/core/Typography'
import { makeStyles } from '@material-ui/core/styles'

import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
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

export default function Version() {
  const classes = useStyles()

  return (
    <div className="container">
      <Head>
        <title>Rules</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom style={{ marginBottom: 26 }}>
            Rules
            <Typography variant="caption" component="sup" style={{ marginLeft: 8, color: 'cyan' }}>
              Sep 19, 2021
            </Typography>
            <Typography>
              To keep the community fair and a little bit safe for both sellers and buyers. Here are
              some written rules to follow.
            </Typography>
          </Typography>

          <div style={{ marginBottom: 26 }}>
            <Typography component="h2" style={{ fontWeight: 'bold' }}>
              Selling Rules
            </Typography>
            <Typography>
              Breaking one of the rules will be punishable by <u>Account Suspension</u> to{' '}
              <u>Permanent Ban</u>.
            </Typography>

            <Typography>
              <ol>
                <li>No misleading post or abusive behavior.</li>
                <li>No multiple reservations on a single item or stock.</li>
                <li>No reservation cancellation without prior notice.</li>
              </ol>
            </Typography>
          </div>

          <div>
            <Typography component="h2" style={{ fontWeight: 'bold' }}>
              Buying Rules
            </Typography>
            <Typography>
              Breaking the rule will reflect a negative impact on buyer&apos;s profile.
            </Typography>

            <Typography>
              <ol>
                <li>No reservation cancellation without prior notice.</li>
              </ol>
            </Typography>
          </div>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
