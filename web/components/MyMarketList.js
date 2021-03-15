import React, { useContext } from 'react'
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
import Snackbar from '@material-ui/core/Snackbar'
import Alert from '@material-ui/lab/Alert'
import * as format from '@/lib/format'
import Button from '@/components/Button'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'
import MarketUpdateDialog from '@/components/MarketUpdateDialog'
import TableSearchInput from '@/components/TableSearchInput'
import Link from '@/components/Link'
import AppContext from '@/components/AppContext'
import {
  VERIFIED_DELIVERY_MAP_TEXT,
  VERIFIED_INVENTORY_MAP_ICON,
  VERIFIED_INVENTORY_MAP_TEXT,
} from '@/constants/verified'

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

export default function MyMarketList({ datatable, loading, error, onSearchInput, onReload }) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const [currentMarket, setCurrentMarket] = React.useState(null)
  const [notifOpen, setNotifOpen] = React.useState(false)

  const handleUpdateClick = marketIdx => {
    setCurrentMarket(datatable.data[marketIdx])
  }
  const handleUpdateSuccess = () => {
    setNotifOpen(true)
    onReload()
  }

  const handleNotifClose = () => {
    setNotifOpen(false)
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
                  loading={loading}
                  onInput={onSearchInput}
                  color="secondary"
                  placeholder="Filter active listings"
                />
              </TableHeadCell>
              {!isMobile && (
                <>
                  <TableHeadCell align="right">Listed</TableHeadCell>
                  <TableHeadCell align="right">Price</TableHeadCell>
                  <TableHeadCell align="center" width={70} />
                </>
              )}
            </TableRow>
          </TableHead>
          <TableBody style={loading ? { opacity: 0.5 } : null}>
            {error && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  Error retrieving data
                  <br />
                  <Typography variant="caption" color="textSecondary">
                    {format.errorSimple(error)}
                  </Typography>
                </TableCell>
              </TableRow>
            )}

            {!error && datatable.data.length === 0 && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  No Result
                </TableCell>
              </TableRow>
            )}

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
                      image={market.item.image}
                      width={77}
                      height={55}
                      title={market.item.name}
                      rarity={market.item.rarity}
                    />
                    <div>
                      <strong>{market.item.name}</strong>
                      <span>{VERIFIED_INVENTORY_MAP_ICON[market.inventory_status]}</span>
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
                          {format.dateFromNow(market.created_at)}
                        </Typography>
                      </TableCell>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {format.amount(market.price, market.currency)}
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
                      <Typography variant="caption" color="textSecondary">
                        <u>Update</u>
                      </Typography>
                    </TableCell>
                  )}
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>
      <MarketUpdateDialog
        open={!!currentMarket}
        market={currentMarket}
        onClose={() => handleUpdateClick(null)}
        onRemove={() => onReload()}
        onSuccess={handleUpdateSuccess}
      />
      <Snackbar open={notifOpen} autoHideDuration={6000} onClose={handleNotifClose}>
        <Alert onClose={handleNotifClose} variant="filled" severity="success">
          Item updated successfully! Check your{' '}
          <Link style={{ textDecoration: 'underline' }} href="/my-listings#reserved">
            Reserved Items
          </Link>
        </Alert>
      </Snackbar>
    </>
  )
}
MyMarketList.propTypes = {
  datatable: PropTypes.object.isRequired,
  onSearchInput: PropTypes.func,
  onReload: PropTypes.func,
  loading: PropTypes.bool,
  error: PropTypes.string,
}
MyMarketList.defaultProps = {
  onSearchInput: () => {},
  onReload: () => {},
  loading: false,
  error: null,
}
