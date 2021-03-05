import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import useSWR from 'swr'
import { fetcher, REPORTS } from '@/service/api'
import { dateFromNow } from '@/lib/format'

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

const filter = {
  sort: 'created_at:desc',
}

export default function Feedback() {
  const classes = useStyles()

  const { data: reports, error } = useSWR([REPORTS, filter], fetcher)

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Feedback Results
          </Typography>

          {error && <Typography color="error">{error}</Typography>}

          <ol>
            {reports &&
              reports.data &&
              reports.data.map(report => (
                <li key={report.id}>
                  <Typography>
                    {report.text} <sup>{dateFromNow(report.created_at)}</sup>
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
