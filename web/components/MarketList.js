import React from 'react'
import PropTypes from 'prop-types'
import useSWR from 'swr'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import Avatar from '@material-ui/core/Avatar'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import { CDN_URL, MARKETS, fetcher } from '@/service/api'
import Link from '@/components/Link'
import BuyButton from '@/components/BuyButton'
import TableHeadCell from '@/components/TableHeadCell'
import { MARKET_STATUS_LIVE } from '../constants/market'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'flex',
    padding: theme.spacing(2),
  },
  avatar: {
    marginRight: theme.spacing(1.5),
  },
}))

const marketFilter = { sort: 'price', status: MARKET_STATUS_LIVE }

export default function MarketList({ itemID = '' }) {
  const classes = useStyles()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))

  marketFilter.item_id = itemID
  const { data: listings, error } = useSWR([MARKETS, marketFilter], (u, f) => fetcher(u, f))

  if (error) {
    return <p>Error</p>
  }

  if (!listings) {
    return <p>Loading...</p>
  }

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableHeadCell>Seller</TableHeadCell>
            <TableHeadCell align="right">Price</TableHeadCell>
            <TableHeadCell align="right" width={156} />
          </TableRow>
        </TableHead>
        <TableBody>
          {listings.data.map(market => (
            <TableRow key={market.id} hover>
              <TableCell component="th" scope="row" padding="none">
                <Link href="/user/[id]" as={`/user/${market.user.steam_id}`} disableUnderline>
                  <div className={classes.seller}>
                    {!isMobile && (
                      <Avatar className={classes.avatar} src={CDN_URL + market.user.avatar} />
                    )}
                    <div>
                      <strong>{market.user.name}</strong>
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {market.user.steam_id}
                      </Typography>
                    </div>
                  </div>
                </Link>
              </TableCell>
              <TableCell align="right">
                <Typography variant="body2">${market.price.toFixed(2)}</Typography>
              </TableCell>
              <TableCell align="right">
                <BuyButton variant="contained">Contact Seller</BuyButton>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
MarketList.propTypes = {
  itemID: PropTypes.string.isRequired,
}
