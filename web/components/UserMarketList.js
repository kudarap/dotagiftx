import React from 'react'
import PropTypes from 'prop-types'
import useSWR from 'swr'
import { makeStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import { MARKETS, fetcher } from '@/service/api'
import Link from '@/components/Link'
import BuyButton from '@/components/BuyButton'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import { MARKET_STATUS_LIVE } from '../constants/market'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'inline-flex',
  },
  link: {
    padding: theme.spacing(2),
  },
}))

const marketFilter = {
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
}

export default function UserMarketList({ userID = '' }) {
  const classes = useStyles()

  marketFilter.user_id = userID
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
            <TableHeadCell>Sell Listings ({listings.total_count})</TableHeadCell>
            <TableHeadCell align="right">Price</TableHeadCell>
            <TableHeadCell align="right" width={156} />
          </TableRow>
        </TableHead>
        <TableBody>
          {listings.data.map(market => (
            <TableRow key={market.id} hover>
              <TableCell component="th" scope="row" padding="none">
                <Link href="/item/[slug]" as={`/item/${market.item.slug}`} disableUnderline>
                  <div className={classes.link}>
                    <strong>{market.item.name}</strong>
                    <br />
                    <Typography variant="caption" color="textSecondary">
                      {market.item.hero}
                    </Typography>
                    <RarityTag rarity={market.item.rarity} />
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
UserMarketList.propTypes = {
  userID: PropTypes.string.isRequired,
}
