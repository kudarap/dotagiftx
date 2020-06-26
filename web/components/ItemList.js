import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles, withStyles } from '@material-ui/core/styles'
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

const useStyles = makeStyles({
  th: {
    cursor: 'pointer',
  },
  pagination: {
    textAlign: 'right',
  },
})

const StyledTableCell = withStyles(theme => ({
  head: {
    textTransform: 'uppercase',
    color: theme.palette.text.secondary,
  },
}))(TableCell)

export default function ItemList({ items = [] }) {
  const classes = useStyles()

  return (
    <>
      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="items table">
          <TableHead>
            <TableRow>
              <StyledTableCell>Name</StyledTableCell>
              <StyledTableCell align="right">Qty</StyledTableCell>
              <StyledTableCell align="right">Price</StyledTableCell>
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
                <TableCell align="right">
                  <Typography variant="body2" color="textSecondary">
                    {item.quantity}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Typography variant="body2">${item.lowest_ask}</Typography>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  )
}
ItemList.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object).isRequired,
}
