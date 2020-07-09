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
import { CDN_URL } from '@/service/api'

const useStyles = makeStyles(theme => ({
  th: {
    cursor: 'pointer',
  },
  pagination: {
    textAlign: 'right',
  },
  link: {
    padding: theme.spacing(2, 2, 2, 0),
    display: 'flex',
  },
  image: {
    margin: theme.spacing(-1, 1, -1, 0),
    width: 77,
    display: 'flex',
    lineHeight: 1,
    flexShrink: 0,
    overflow: 'hidden',
    userSelect: 'none',
  },
}))

function ItemImage({ image, title }) {
  const classes = useStyles()

  const imgStyle = {
    color: 'transparent',
    width: '100%',
    height: '100%',
    objectFit: 'cover',
    textAlign: 'center',
    textIndent: '10000px',
  }

  return (
    <div className={classes.image}>
      <img src={`${CDN_URL + image}/200x100`} alt={title} style={imgStyle} />
    </div>
  )
}

export default function CatalogList({ items = [], variant }) {
  const classes = useStyles()

  const isRecentMode = variant === 'recent'

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="items table">
        <TableHead>
          <TableRow>
            <TableHeadCell>Item</TableHeadCell>
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
                    <ItemImage image={item.image} title={item.name} />
                    <div>
                      <strong>{item.name}</strong>
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {item.hero}
                      </Typography>
                      <RarityTag rarity={item.rarity} />
                    </div>
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
