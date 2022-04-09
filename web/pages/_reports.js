import React from 'react'
import useSWR from 'swr'
import map from 'lodash/map'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import { fetcher, REPORTS } from '@/service/api'
import { dateFromNow } from '@/lib/format'
import { REPORT_TYPE_MAP_TEXT, REPORT_TYPE_SURVEY } from '@/constants/report'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
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
  const { classes } = useStyles()

  const { data: reports, error } = useSWR([REPORTS, filter], fetcher)

  // tally report data base on text
  if (reports && reports.data) {
    reports.data.forEach(report => {
      if (report.type !== REPORT_TYPE_SURVEY) {
        return
      }

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
                    from <Link href={`/profiles/${report.user.steam_id}`}>{report.user.name}</Link>
                    &nbsp;
                    {dateFromNow(report.created_at)}
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
