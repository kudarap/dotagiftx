import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
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
import { amount } from '@/lib/format'
import Button from '@/components/Button'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'
import ReserveUpdateDialog from '@/components/ReserveUpdateDialog'
import TableSearchInput from '@/components/TableSearchInput'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'inline-flex',
  },
  item: {
    padding: theme.spacing(2, 2, 2, 0),
    display: 'flex',
    cursor: 'pointer',
  },
  image: {
    margin: theme.spacing(-1, 1, -1, 1),
    width: 77,
    height: 55,
  },
}))

export default function ReservationList({ datatable, loading, error, onSearchInput }) {
  const classes = useStyles()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))

  const [currentMarket, setCurrentMarket] = React.useState(null)
  const [notifOpen, setNotifOpen] = React.useState(false)

  const handleUpdateClick = marketIdx => {
    setCurrentMarket(datatable.data[marketIdx])
  }

  const handleNotifClose = () => {
    setNotifOpen(false)
  }

  const handleUpdateSuccess = () => {
    setNotifOpen(true)
    onSearchInput('')
  }

  return (
    <>
      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableHeadCell padding="none" colSpan={isMobile ? 2 : 1}>
                <TableSearchInput
                  fullWidth
                  onInput={onSearchInput}
                  color="secondary"
                  placeholder="Filter items"
                />
              </TableHeadCell>
              {!isMobile && (
                <>
                  <TableHeadCell align="right">Reserved</TableHeadCell>
                  <TableHeadCell align="right">Price</TableHeadCell>
                  <TableHeadCell align="center" width={70} />
                </>
              )}
            </TableRow>
          </TableHead>
          <TableBody>
            {datatable.data &&
              datatable.data.map((market, idx) => (
                <TableRow key={market.id} hover>
                  <TableCell
                    component="th"
                    scope="row"
                    padding="none"
                    className={classes.item}
                    onClick={() => handleUpdateClick(idx)}>
                    <ItemImage
                      className={classes.image}
                      image={`/200x100/${market.item.image}`}
                      title={market.item.name}
                      rarity={market.item.rarity}
                    />
                    <div>
                      <strong>{market.item.name}</strong>
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {market.item.hero}
                      </Typography>
                      <RarityTag rarity={market.item.rarity} />
                    </div>
                  </TableCell>
                  {!isMobile ? (
                    <>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {moment(market.updated_at).fromNow()}
                        </Typography>
                      </TableCell>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {amount(market.price, market.currency)}
                        </Typography>
                      </TableCell>
                      <TableCell align="center">
                        <Button variant="outlined" onClick={() => handleUpdateClick(idx)}>
                          Update
                        </Button>
                      </TableCell>
                    </>
                  ) : (
                    <TableCell
                      align="right"
                      onClick={() => handleUpdateClick(idx)}
                      style={{ cursor: 'pointer' }}>
                      <Typography variant="body2" color="secondary">
                        {format.amount(market.price, market.currency)}
                      </Typography>
                      <Typography variant="caption" color="textSecondary" noWrap>
                        {moment(market.updated_at).fromNow()}
                      </Typography>
                    </TableCell>
                  )}
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>
      <ReserveUpdateDialog
        open={!!currentMarket}
        market={currentMarket}
        onClose={() => handleUpdateClick(null)}
      />
    </>
  )
}
ReservationList.propTypes = {
  datatable: PropTypes.object.isRequired,
  error: PropTypes.string,
}
ReservationList.defaultProps = {
  error: null,
}
