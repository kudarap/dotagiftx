import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import Paper from '@mui/material/Paper'
import Alert from '@mui/material/Alert'
import Skeleton from '@mui/material/Skeleton'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import InternalUserCard from '@/components/InternalUserCard'
import { user } from '@/service/api'
import { dateCalendar, dateTime } from '@/lib/format'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  list: {
    listStyle: 'none',
    '& li:before': {
      content: `'🛡️ '`,
    },
    paddingLeft: theme.spacing(3),
    marginTop: 0,
  },
}))

function createRate(payment, serviceFee, minimumFee, pulloutFee, disputeFee) {
  return { payment, serviceFee, minimumFee, pulloutFee, disputeFee }
}

// paypal fees 4.4% + 0.30
const tableRatesRev0 = [
  createRate('PayPal', '+10%', '$1.00', '10% + 4.4% + $0.30', '4.4% + $0.30'),
  createRate('Mann Co. Supply Crate Key (TF key)', '+15%', '1 Key', '15%', 'None'),
  createRate('Crypto', 'TBD', 'TBD', 'TBD', 'TBD'),
]

const tableRates = [
  createRate('PayPal', '15% + $0.60', 'None', 'None', 'None'),
  createRate('Mann Co. Supply Crate Key (TF key)', '10%', 'None', 'None', 'None'),
  createRate('Crypto', '5%', 'None', 'None', 'None'),
]

const middlemanUserIds = ['76561198078354099']

const updatedAt = new Date('2026-08-08')

export default function Middleman() {
  const { classes } = useStyles()

  const [middleman, setMiddleman] = React.useState([])
  const [loading, setLoading] = React.useState(true)

  React.useEffect(() => {
    async function fetchUser(id) {
      return new Promise(resolve => {
        user(id).then(u => {
          resolve({
            id: u.steam_id,
            name: u.name,
            img: u.avatar,
            boons: ['MIDDLEMAN_TAG'],
            discordURL: 'https://discord.gg/b79zMpjjc5',
            createdAt: u.created_at,
          })
        })
      })
    }

    async function fetchUsers(array) {
      const res = await Promise.all(
        array.map(async item => {
          const v = await fetchUser(item)
          return v
        })
      )
      return res
    }

    fetchUsers(middlemanUserIds).then(res => {
      setMiddleman(res)
      setLoading(false)
    })
  }, [])

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>Middleman :: {APP_NAME}</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Middleman
            <Typography variant="body2" color="textSecondary">
              {dateCalendar(updatedAt)}
            </Typography>
          </Typography>
          <Typography gutterBottom>
            The profile listed below is the only official middleman service provider of the site.
            Please read the terms of this service carefully.
          </Typography>

          {loading &&
            middlemanUserIds.map(id => (
              <Skeleton key={id}>
                <InternalUserCard boons={['MIDDLEMAN_TAG']} />
              </Skeleton>
            ))}

          {middleman.map(row => (
            <InternalUserCard key={row.id} {...row} />
          ))}

          <Typography component="h2" variant="h6" gutterBottom>
            Service rates
            <Typography variant="body2" color="orange">
              {dateTime(updatedAt)}
            </Typography>
          </Typography>
          <Alert severity="warning">
            Rates are subject to change without prior notice, outstanding transaction fees will
            remain as it is.
          </Alert>

          <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
              <TableHead>
                <TableRow>
                  <TableCell>Payment type</TableCell>
                  <TableCell align="right">Service fee</TableCell>
                  <TableCell align="center">Minimum fee</TableCell>
                  <TableCell align="center">Cancel penalty</TableCell>
                  <TableCell align="center">Abandon penalty</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {tableRates.map(row => (
                  <TableRow
                    key={row.payment}
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                    <TableCell component="th" scope="row">
                      {row.payment}
                    </TableCell>
                    <TableCell align="right">{row.serviceFee}</TableCell>
                    <TableCell align="center">{row.minimumFee}</TableCell>
                    <TableCell align="center">{row.pulloutFee}</TableCell>
                    <TableCell align="center">{row.disputeFee}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
          <br />

          {/* <Typography component="h2" variant="h6">
            Calculator
          </Typography>
          <Typography color="textSecondary">
            Market place for Dota 2 Giftables, items that can only be gift or gift-once are probably
            belong here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread
            with a kick.
          </Typography>
          <br />*/}

          <Typography component="h2" variant="h6">
            How to proceed?
          </Typography>
          <Typography color="textSecondary">
            The middleman service is recommended if you want to pay a seller after the 30 day friend
            requirement has been completed and both parties are available to complete the trade.
            Simply join the Discord server and request a middleman to ensure a secure transaction.
          </Typography>
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
