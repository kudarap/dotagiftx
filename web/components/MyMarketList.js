import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
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
import Button from '@/components/Button'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'inline-flex',
  },
  link: {
    [theme.breakpoints.down('xs')]: {
      paddingLeft: theme.spacing(2),
    },
    padding: theme.spacing(2, 2, 2, 0),
    display: 'flex',
  },
  image: {
    margin: theme.spacing(-1, 1, -1, 1),
    width: 77,
  },
}))

export default function MyMarketList({ datatable, error }) {
  const classes = useStyles()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))

  if (error) {
    return <p>Error</p>
  }

  if (!datatable) {
    return <p>Loading...</p>
  }

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableHeadCell>Sell Listings ({datatable.total_count})</TableHeadCell>
            <TableHeadCell align="right">Listed</TableHeadCell>
            <TableHeadCell align="right">Price</TableHeadCell>
            <TableHeadCell align="right" width={156} />
          </TableRow>
        </TableHead>
        <TableBody>
          {datatable.data.map(market => (
            <TableRow key={market.id} hover>
              <TableCell component="th" scope="row" padding="none">
                <Link
                  className={classes.link}
                  href="/item/[slug]"
                  as={`/item/${market.item.slug}`}
                  disableUnderline>
                  {!isMobile && (
                    <ItemImage
                      className={classes.image}
                      image={`/200x100/${market.item.image}`}
                      title={market.item.name}
                      rarity={market.item.rarity}
                    />
                  )}
                  <div>
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
                <Typography variant="body2">{format.dateFromNow(market.created_at)}</Typography>
              </TableCell>
              <TableCell align="right">
                <Typography variant="body2">${market.price.toFixed(2)}</Typography>
              </TableCell>
              <TableCell align="right">
                <Button variant="contained">Edit</Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
MyMarketList.propTypes = {
  datatable: PropTypes.object.isRequired,
  error: PropTypes.string,
}
MyMarketList.defaultProps = {
  error: null,
}
