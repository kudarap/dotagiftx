import React, { useContext } from 'react'
import moment from 'moment'
import PropTypes from 'prop-types'
import bidColor from '@material-ui/core/colors/teal'
import { makeStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import * as format from '@/lib/format'
import Link from '@/components/Link'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'
import AppContext from '@/components/AppContext'

const useStyles = makeStyles(theme => ({
  th: {
    cursor: 'pointer',
  },
  pagination: {
    textAlign: 'right',
  },
  link: {
    [theme.breakpoints.down('xs')]: {
      paddingLeft: theme.spacing(0),
    },
    padding: theme.spacing(2, 2, 2, 0),
    display: 'flex',
  },
  image: {
    margin: theme.spacing(-1, 1, -1, 1),
    width: 77,
    height: 55,
  },
}))

export default function CatalogList({ items = [], loading, error, variant, bidType }) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const isRecentMode = variant === 'recent'

  const itemURLSuffix = bidType ? '/buyorders' : ''

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="items table">
        <TableHead>
          <TableRow>
            <TableHeadCell>Item</TableHeadCell>
            {!isMobile && (
              <TableHeadCell align="right">
                {/* eslint-disable-next-line no-nested-ternary */}
                {isRecentMode ? (bidType ? 'Ordered' : 'Listed') : 'Qty'}
              </TableHeadCell>
            )}
            <TableHeadCell align="right">{bidType ? 'Buy Price' : 'Price'}</TableHeadCell>
          </TableRow>
        </TableHead>
        <TableBody style={loading ? { opacity: 0.5 } : null}>
          {error && (
            <TableRow>
              <TableCell align="center" colSpan={3}>
                Error retrieving data
                <br />
                <Typography variant="caption" color="textSecondary">
                  {error}
                </Typography>
              </TableCell>
            </TableRow>
          )}

          {!error && items && items.length === 0 && (
            <TableRow>
              <TableCell align="center" colSpan={3}>
                No Result
              </TableCell>
            </TableRow>
          )}

          {items &&
            items.map(item => (
              <TableRow key={item.id} hover>
                <TableCell className={classes.th} component="th" scope="row" padding="none">
                  <Link
                    className={classes.link}
                    href={`/[slug]${itemURLSuffix}`}
                    as={`/${item.slug}${itemURLSuffix}`}
                    disableUnderline>
                    <ItemImage
                      className={classes.image}
                      image={item.image}
                      width={77}
                      height={55}
                      title={item.name}
                      rarity={item.rarity}
                    />

                    <div>
                      <strong>{item.name}</strong>
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {item.hero}
                      </Typography>
                      <RarityTag rarity={item.rarity} />
                    </div>
                  </Link>
                </TableCell>

                {!isMobile && (
                  <TableCell align="right">
                    <Typography variant="body2" color="textSecondary">
                      {/* eslint-disable-next-line no-nested-ternary */}
                      {isRecentMode
                        ? moment(bidType ? item.recent_bid : item.recent_ask).fromNow()
                        : bidType
                        ? item.bid_count
                        : item.quantity}
                    </Typography>
                  </TableCell>
                )}

                <TableCell align="right">
                  <Typography variant="body2" style={bidType ? { color: bidColor.A200 } : null}>
                    {format.amount(bidType ? item.highest_bid : item.lowest_ask, 'USD')}
                  </Typography>
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
  loading: PropTypes.bool,
  error: PropTypes.string,
  bidType: PropTypes.bool,
}
CatalogList.defaultProps = {
  variant: '',
  loading: false,
  error: null,
  bidType: false,
}
