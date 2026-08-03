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
import Skeleton from '@mui/material/Skeleton'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import InternalUserCard from '@/components/InternalUserCard'
import type { InternalUserCardProps } from '@/components/InternalUserCard'
import { user } from '@/service/api'
import type { Profile } from '@/lib/types'

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
  payment: string,
  serviceFee: string,
  minimumFee: string,
  pulloutFee: string,
  disputeFee: string
) {
  return { payment, serviceFee, minimumFee, pulloutFee, disputeFee }
}

// paypal fees 4.4% + 0.30
const tableRates = [
  createRate('PayPal', '+10%', '$1.00', '10% + 4.4% + $0.30', '4.4% + $0.30'),
  createRate('Mann Co. Supply Crate Key (TF key)', '+15%', '1 Key', '15%', 'None'),
  createRate('Crypto', 'TBD', 'TBD', 'TBD', 'TBD'),
]

const middlemanUserIds = ['76561198088587178']

export default function Middleman() {
  const { classes } = useStyles()

  const [middleman, setMiddleman] = React.useState<InternalUserCardProps[]>([])
  const [loading, setLoading] = React.useState(true)

  React.useEffect(() => {
    async function fetchUser(id: string) {
      return new Promise(resolve => {
        user(id).then(u => {
          const profile = u as Profile
          resolve({
            id: profile.steam_id,
            name: profile.name,
            img: profile.avatar,
            boons: ['MIDDLEMAN_TAG'],
            discordURL: 'https://discord.gg/b79zMpjjc5',
            createdAt: profile.created_at,
          })
        })
      })
    }

    async function fetchUsers(array: string[]) {
      const res = await Promise.all(
        array.map(async item => {
          const v = await fetchUser(item)
          return v
        })
      )
      return res
    }

    fetchUsers(middlemanUserIds).then(res => {
      setMiddleman(res as InternalUserCardProps[])
      setLoading(false)
    })
  }, [])

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

          {loading &&
            middlemanUserIds.map(id => (
              <Skeleton key={id}>
                <InternalUserCard boons={['MIDDLEMAN_TAG']} />
              </Skeleton>
            ))}

          {middleman.map(row => (
            <InternalUserCard key={row.id} {...row} />
          ))}

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
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                  >
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
        </Container>
      </main>

      <Footer />
    </>
  )
}
