import React from 'react'
import Head from 'next/head'
import PropTypes from 'prop-types'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import Link from '@mui/material/Link'
import Paper from '@mui/material/Paper'
import Alert from '@mui/material/Alert'
import Skeleton from '@mui/material/Skeleton'
import Grid from '@mui/material/Grid'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import Box from '@mui/material/Box'
import TextField from '@mui/material/TextField'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import InternalUserCard from '@/components/InternalUserCard'
import Button from '@/components/Button'
import { user } from '@/service/api'
import { amount, dateCalendar, dateTimeFull } from '@/lib/format'

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

function createRate(
  payment,
  serviceFee,
  minimumFee,
  pulloutFee,
  disputeFee,
  percent = 0,
  flatFee = 0
) {
  return { payment, serviceFee, minimumFee, pulloutFee, disputeFee, percent, flatFee }
}

// paypal fees 4.4% + 0.30
const tableRatesRev0 = [
  createRate('PayPal', '+10%', '$1.00', '10% + 4.4% + $0.30', '4.4% + $0.30'),
  createRate('Mann Co. Supply Crate Key (TF key)', '+15%', '1 Key', '15%', 'None'),
  createRate('Crypto', 'TBD', 'TBD', 'TBD', 'TBD'),
]

const tableRates = [
  createRate('PayPal', '15% + $0.60', 'None', 'None', 'None'),
  createRate('TF keys, Dota 2 and Rust items', '10%', 'None', 'None', 'None'),
  createRate('Crypto', '5%', 'None', 'None', 'None'),
]

const computeRates = [
  {
    id: 'paypal',
    label: 'PayPal 15% + $0.60',
    percetage: 0.15,
    flat: 0.6,
  },
  {
    id: 'steamItems',
    label: 'TF keys, Dota 2 and Rust items 10%',
    percetage: 0.1,
    flat: 0,
  },
  {
    id: 'crypto',
    label: 'Crypto 5%',
    percetage: 0.05,
    flat: 0,
  },
]

const middlemanUserIds = ['76561198078354099']
const updatedAt = new Date('2026-08-08')
const middlemanDiscordURL = 'https://discord.gg/b79zMpjjc5'

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
            discordURL: middlemanDiscordURL,
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
              {dateTimeFull(updatedAt)}
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

          <Calculator rates={computeRates} />

          <Typography component="h2" variant="h6">
            How to proceed?
          </Typography>
          <Typography color="textSecondary">
            The middleman service is recommended if you want to pay a seller after the 30 day friend
            requirement has been completed and both parties are available to complete the trade.
            Simply join the Discord server and request a{' '}
            <strong style={{ color: '#2ecc71', fontFamily: 'monospace' }}>@middleman</strong> to
            ensure a secure transaction.
          </Typography>
          <br />

          <Box>
            <Button
              component={Link}
              target="_blank"
              rel="noreferrer noopener"
              href={middlemanDiscordURL}
              size="large"
              variant="outlined"
              color="secondary">
              Continue To Discord
            </Button>
          </Box>
        </Container>
      </main>

      <Footer />
    </>
  )
}

function Calculator({ rates }) {
  const [payment, setPayment] = React.useState(rates[0] ? rates[0].id : '')
  const [price, setPrice] = React.useState(1)

  const [computed, setComputed] = React.useState({ fee: 0, total: 0 })
  React.useEffect(() => {
    const rate = rates.find(rate => rate.id == payment)
    const fee = Number(rate.flat + price * rate.percetage)
    const total = Number(price) + fee
    setComputed({ fee, total })
  }, [payment, price, rates])

  return (
    <>
      <Typography component="h2" variant="h6" gutterBottom>
        Calculator
      </Typography>
      <Alert severity="info">
        Estimate the service fee for your transaction. Fee is computed based on the rates above.
      </Alert>
      <br />
      <Grid container spacing={2}>
        <Grid item xs={12} sm={6}>
          <FormControl fullWidth color="secondary" variant="outlined">
            <InputLabel id="calculator-payment-type-label">Payment type</InputLabel>
            <Select
              labelId="calculator-payment-type-label"
              id="calculator-payment-type"
              variant="outlined"
              value={payment}
              onChange={e => setPayment(e.target.value)}
              label="Payment type">
              {rates.map(row => (
                <MenuItem key={row.id} value={row.id}>
                  {row.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sm={6}>
          <TextField
            fullWidth
            color="secondary"
            variant="outlined"
            label="Price"
            type="number"
            inputProps={{ min: 0, step: 'any' }}
            value={price}
            onChange={e => setPrice(e.target.value)}
          />
        </Grid>
      </Grid>
      <br />
      <Paper
        sx={{ padding: 2, display: 'flex', justifyContent: 'space-around', textAlign: 'center' }}>
        <div>
          <Typography variant="body2" color="textSecondary">
            Service fee
          </Typography>
          <Typography variant="h6">{amount(computed.fee, 'USD')}</Typography>
        </div>
        <div>
          <Typography variant="body2" color="textSecondary">
            Total amount
          </Typography>
          <Typography variant="h6">{amount(computed.total, 'USD')}</Typography>
        </div>
      </Paper>
      <br />
    </>
  )
}
Calculator.propTypes = {
  rates: PropTypes.array,
}
Calculator.defaultProps = {
  rates: [],
}
