import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { makeStyles } from '@material-ui/core/styles'
import { debounce } from '@material-ui/core'
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
import {
  VERIFIED_INVENTORY_MAP_ICON,
  VERIFIED_INVENTORY_VERIFIED_RESELL,
} from '@/constants/verified'
import * as format from '@/lib/format'
import { amount } from '@/lib/format'
import Button from '@/components/Button'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'
import ReserveUpdateDialog from '@/components/ReserveUpdateDialog'
import TableSearchInput from '@/components/TableSearchInput'
import Link from '@/components/Link'
import AppContext from '@/components/AppContext'
import { VerifiedStatusPopover } from '@/components/VerifiedStatusCard'

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

export default function ReservationList({ datatable, loading, error, onSearchInput, onReload }) {
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

  const [currentIndex, setIndex] = React.useState(null)
  const [anchorEl, setAnchorEl] = React.useState(null)
  const debouncePopoverClose = debounce(() => {
    setAnchorEl(null)
    setIndex(null)
  }, 150)
  const handlePopoverOpen = event => {
    debouncePopoverClose.clear()
    setIndex(Number(event.currentTarget.dataset.index))
    setAnchorEl(event.currentTarget)
  }
  const handlePopoverClose = () => {
    setAnchorEl(null)
    setIndex(null)
  }
  const open = Boolean(anchorEl)
  const popoverElementID = open ? 'verified-status-popover' : undefined

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
                  placeholder="Filter reserved items"
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
                      <span
                        aria-owns={popoverElementID}
                        aria-haspopup="true"
                        data-index={idx}
                        onMouseLeave={debouncePopoverClose}
                        onMouseEnter={handlePopoverOpen}>
                        {market.resell
                          ? VERIFIED_INVENTORY_MAP_ICON[VERIFIED_INVENTORY_VERIFIED_RESELL]
                          : VERIFIED_INVENTORY_MAP_ICON[market.inventory_status]}
                      </span>

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
        onCancel={() => onReload()}
        onSuccess={handleUpdateSuccess}
      />

      <VerifiedStatusPopover
        id={popoverElementID}
        open={open}
        anchorEl={anchorEl}
        onClose={handlePopoverClose}
        onMouseEnter={() => debouncePopoverClose.clear()}
        market={datatable.data[currentIndex]}
      />

      <Snackbar open={notifOpen} autoHideDuration={6000} onClose={handleNotifClose}>
        <Alert onClose={handleNotifClose} variant="filled" severity="success">
          Item updated successfully! Check your{' '}
          <Link style={{ textDecoration: 'underline' }} href="/my-listings#delivered">
            History Items
          </Link>
        </Alert>
      </Snackbar>
    </>
  )
}
ReservationList.propTypes = {
  datatable: PropTypes.object.isRequired,
  onSearchInput: PropTypes.func,
  onReload: PropTypes.func,
  loading: PropTypes.bool,
  error: PropTypes.string,
}
ReservationList.defaultProps = {
  onSearchInput: () => {},
  onReload: () => {},
  loading: false,
  error: null,
}
