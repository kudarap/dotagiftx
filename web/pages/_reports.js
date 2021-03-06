import React from 'react'
import useSWR from 'swr'
import map from 'lodash/map'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { fetcher, REPORTS } from '@/service/api'
import { dateFromNow } from '@/lib/format'
import { REPORT_TYPE_MAP_TEXT } from '@/constants/report'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
    // background: 'url("/icon.png") no-repeat bottom right',
    // backgroundSize: 100,
  },
}))

const tallyVotes = {}

const filter = {
  sort: 'created_at:desc',
  limit: 100,
}

export default function Feedback() {
  const classes = useStyles()

  const { data: reports, error } = useSWR([REPORTS, filter], fetcher)

  // tally report data base on text
  if (reports && reports.data) {
    reports.data.forEach(report => {
      if (!tallyVotes[report.text]) {
        tallyVotes[report.text] = 1
        return
      }

      tallyVotes[report.text]++
    })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Feedback Results
          </Typography>

          {error && <Typography color="error">{error}</Typography>}

          {reports &&
            reports.data &&
            map(tallyVotes, (text, score) => {
              return (
                <Typography color="secondary">
                  {text}x {score}
                </Typography>
              )
            })}

          <ol>
            {reports &&
              reports.data &&
              reports.data.map(report => (
                <li key={report.id}>
                  <Typography color="textSecondary">
                    <Typography color="textPrimary" component="span">
                      {REPORT_TYPE_MAP_TEXT[report.type].toUpperCase()}
                    </Typography>
                    {` ${report.text} `}
                    <em>
                      <small>{dateFromNow(report.created_at)}</small>
                    </em>
                  </Typography>
                </li>
              ))}
          </ol>
        </Container>
      </main>

      <Footer />
    </>
  )
}
