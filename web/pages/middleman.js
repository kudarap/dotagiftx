import React from 'react'
import PropTypes from 'prop-types'
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
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import InternalUserCard from '@/components/InternalUserCard'
import { user } from '@/service/api'

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

const middlemanUserIds = [
  // '76561198287849998',
]

// paypal fees 4.4% + 0.30
const tableRates = [
  createRate('PayPal', '+10%', '$1.00', '10% + 4.4% + $0.30', '4.4% + $0.30'),
  createRate('Mann Co. Supply Crate Key (TF key)', '+15%', '1 Key', '15%', 'None'),
  createRate('Crypto', 'TBD', 'TBD', 'TBD', 'TBD'),
]

export default function Middleman({ users }) {
  const { classes } = useStyles()

  const middleman = users.map(row => ({
    id: row.steam_id,
    name: row.name,
    img: row.avatar,
    boons: row.boons,
    discordURL: 'https://discord.gg/b79zMpjjc5',
  }))

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Middleman</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom>
            Middleman
            <Typography variant="body2" color="textSecondary">
              Updated June 25, 2025
            </Typography>
          </Typography>
          <Typography color="textSecondary">
            The profile listed below is the only official middleman service provider of the site.
            Please read the terms of this service carefully.
          </Typography>
          <br />

          {middleman.map(row => (
            <InternalUserCard key={row.id} {...row} />
          ))}
          <InternalUserCard
            id="76561198088587178"
            name="kudarap"
            img="7055bd1d085fdf1ff9e9928df06ec64c1d04c124.jpg"
            boons={['MIDDLEMAN_TAG']}
            discordURL="https://discord.gg/b79zMpjjc5"
          />
          <br />

          {/* <Typography component="h2" variant="h6">
            Calculator
          </Typography>
          <Typography color="textSecondary">
            Market place for Dota 2 Giftables, items that can only be gift or gift-once are probably
            belong here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread
            with a kick.
          </Typography>
          <br />

          <Typography component="h2" variant="h6">
            Process
          </Typography>
          <Typography color="textSecondary" gutterBottom>
            Market place for Dota 2 Giftables, items that can only be gift or gift-once are probably
            belong here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread
            with a kick.
          </Typography>
          <Typography color="textSecondary">
            Market place for Dota 2 Giftables, items that can only be gift or gift-once are probably
            belong here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread
            with a kick.
          </Typography>
          <br /> */}

          <Typography component="h2" variant="h6">
            Service rates
          </Typography>
          <Typography color="textSecondary">
            Rates updated at Sep 17, 2025 and subject to change without prior notice, outstanding
            transaction fees will remain as it is.
          </Typography>
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
            Terms
          </Typography>
          <Typography color="textSecondary">
            Market place for Dota 2 Giftables, items that can only be gift or gift-once are probably
            belong here. If you are on Dota2Trade subreddit, its basically the Giftable Megathread
            with a kick.
          </Typography>
          <br />

          <Typography component="h2" variant="h6">
            FAQs
          </Typography>
          <Box>
            <Typography>What is DotagiftX?</Typography>
            <Typography color="textSecondary" gutterBottom>
              Market place for Dota 2 Giftables, items that can only be gift or gift-once are
              probably belong here. If you are on Dota2Trade subreddit, its basically the Giftable
              Megathread with a kick.
            </Typography>

            <Typography>What is DotagiftX?</Typography>
            <Typography color="textSecondary" gutterBottom>
              Market place for Dota 2 Giftables, items that can only be gift or gift-once are
              probably belong here. If you are on Dota2Trade subreddit, its basically the Giftable
              Megathread with a kick.
            </Typography>
          </Box> */}
        </Container>
      </main>

      <Footer />
    </>
  )
}

Middleman.propTypes = {
  users: PropTypes.arrayOf(PropTypes.object).isRequired,
  error: PropTypes.string,
}
Middleman.defaultProps = {
  users: [],
  error: null,
}

// This gets called on every request
export async function getServerSideProps(context) {
  let users = []
  for (const id of middlemanUserIds) {
    try {
      users.push(await user(id))
    } catch (e) {
      return {
        props: {
          error: e.message,
        },
      }
    }
  }
  return {
    props: {
      users,
    },
  }
}
