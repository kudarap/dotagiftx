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
import TablePagination from '@/components/TableActions'
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

export default function ItemList({
  result = {
    data: [],
    result_count: 0,
    total_count: 0,
  },
  onChangePage,
}) {
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
            {result.data.map(item => (
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
                    {item.name.length}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Typography variant="body2">${item.hero.length.toFixed(2)}</Typography>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        className={classes.pagination}
        colSpan={3}
        count={result.total_count}
        page={1}
        onChangePage={onChangePage}
      />
    </>
  )
}
ItemList.propTypes = {
  result: PropTypes.object.isRequired,
  onChangePage: PropTypes.func,
}
ItemList.defaultProps = {
  onChangePage: () => {},
}
