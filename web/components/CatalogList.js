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

const useStyles = makeStyles(theme => ({
  th: {
    cursor: 'pointer',
  },
  pagination: {
    textAlign: 'right',
  },
  link: {
    padding: theme.spacing(2),
  },
}))

export default function CatalogList({ items = [], variant }) {
  const classes = useStyles()

  const isRecentMode = variant === 'recent'

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="items table">
        <TableHead>
          <TableRow>
            <TableHeadCell>Name</TableHeadCell>
            <TableHeadCell align="right">{isRecentMode ? 'Listed' : 'Qty'}</TableHeadCell>
            <TableHeadCell align="right">Price</TableHeadCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {items.length === 0 && (
            <TableRow>
              <TableCell align="center" colSpan={3}>
                No Result
              </TableCell>
            </TableRow>
          )}

          {items.map(item => (
            <TableRow key={item.id} hover>
              <TableCell className={classes.th} component="th" scope="row" padding="none">
                <Link href="/item/[slug]" as={`/item/${item.slug}`} disableUnderline>
                  <div className={classes.link}>
                    <strong>{item.name}</strong>
                    <br />
                    <Typography variant="caption" color="textSecondary">
                      {item.hero}
                    </Typography>
                    <RarityTag rarity={item.rarity} />
                  </div>
                </Link>
              </TableCell>

              <TableCell align="right">
                <Typography variant="body2" color="textSecondary">
                  {isRecentMode ? moment(item.recent_ask).fromNow() : item.quantity}
                </Typography>
              </TableCell>

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
CatalogList.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object).isRequired,
  variant: PropTypes.string,
}
CatalogList.defaultProps = {
  variant: '',
}
