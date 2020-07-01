import React from 'react'
import moment from 'moment'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import Link from '@/components/Link'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'

const useStyles = makeStyles({
  th: {
    cursor: 'pointer',
  },
  pagination: {
    textAlign: 'right',
  },
})

export default function ItemList({ items = [], variant }) {
  const classes = useStyles()

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="items table">
        <TableHead>
          <TableRow>
            <TableHeadCell>Name</TableHeadCell>
            {variant === 'recent' ? (
              <TableHeadCell>Date</TableHeadCell>
            ) : (
              <TableHeadCell align="right">Qty</TableHeadCell>
            )}
            <TableHeadCell align="right">Price</TableHeadCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {items.map(item => (
            <TableRow key={item.id} hover>
              <TableCell className={classes.th} component="th" scope="row">
                <Link href="/item/[slug]" as={`/item/${item.slug}`} disableUnderline>
                  <>
                    <strong>{item.name}</strong>
                    <br />
                    <Typography variant="caption" color="textSecondary">
                      {item.hero}
                    </Typography>
                    <RarityTag rarity={item.rarity} />
                  </>
                </Link>
              </TableCell>
              {variant === 'recent' ? (
                <TableCell>
                  <Typography variant="body2" color="textSecondary">
                    {moment(item.created_at).fromNow()}
                  </Typography>
                </TableCell>
              ) : (
                <TableCell align="right">
                  <Typography variant="body2" color="textSecondary">
                    {item.quantity}
                  </Typography>
                </TableCell>
              )}

              <TableCell align="right">
                <Typography variant="body2">${item.lowest_ask.toFixed(2)}</Typography>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
ItemList.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object).isRequired,
  variant: PropTypes.string,
}
ItemList.defaultProps = {
  variant: '',
}
