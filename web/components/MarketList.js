import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import { debounce, NoSsr } from '@mui/material'
import { teal as bidColor } from '@mui/material/colors'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import Chip from '@mui/material/Chip'
import {
  VERIFIED_INVENTORY_MAP_ICON,
  VERIFIED_INVENTORY_VERIFIED_RESELL,
} from '@/constants/verified'
import { isDonationGlowExpired, myMarket } from '@/service/api'
import { amount, dateFromNow } from '@/lib/format'
import Link from '@/components/Link'
import Button from '@/components/Button'
import BuyButton from '@/components/BuyButton'
import TableHeadCell from '@/components/TableHeadCell'
import ContactDialog from '@/components/ContactDialog'
import ContactBuyerDialog from '@/components/ContactBuyerDialog'
import { MARKET_STATUS_REMOVED } from '@/constants/market'
import { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import SellButton from '@/components/SellButton'
import { VerifiedStatusPopover } from '@/components/VerifiedStatusCard'
import Avatar from '@/components/Avatar'
import DashTabs from '@/components/DashTabs'
import DashTab from '@/components/DashTab'
import { getUserBadgeFromBoons } from '@/lib/badge'
import SubscriberBadge from './SubscriberBadge'

const useStyles = makeStyles()(theme => ({
  seller: {
    display: 'flex',
    padding: theme.spacing(2),
  },
  avatar: {
    marginRight: theme.spacing(1.5),
  },
  tableHead: {
    // background: '#202a2f',
    background: 'linear-gradient(to right, #9d731f1f, #52c6bb26)',
  },
  tabs: {
    '& .MuiTabs-indicator': {
      background: theme.palette.grey[100],
    },
  },
  tab: {
    width: 168,
    textTransform: 'none',
  },
  sortButtons: {
    display: 'flex',
    '& .MuiChip-root': {
      marginRight: theme.spacing(1),
    },
  },
  activeSortButtons: {
    color: `${theme.palette.grey[800]} !important`,
    background: `${theme.palette.grey[100]} !important`,
  },
}))

export default function MarketList({
  offers,
  buyOrders,
  error,
  loading,
  sort,
  pagination,
  tabIndex,
  onSortChange,
  onTabChange,
}) {
  const { classes } = useStyles()
  const { isMobile, currentAuth } = useContext(AppContext)
  const currentUserID = currentAuth.user_id || null

  const router = useRouter()
  const handleTabChange = (e, value) => {
    onTabChange(value)
  }

  const [currentMarket, setCurrentMarket] = React.useState(null)
  const handleContactClick = marketIdx => {
    let src = offers
    if (tabIndex === 1) {
      src = buyOrders
    }

    setCurrentMarket(src.data[marketIdx])
  }
  const handleRemoveClick = marketIdx => {
    let src = offers
    if (tabIndex === 1) {
      src = buyOrders
    }

    const mktID = src.data[marketIdx].id
    ;(async () => {
      try {
        await myMarket.PATCH(mktID, { status: MARKET_STATUS_REMOVED })
        router.reload()
      } catch (e) {
        console.error(`Error: ${e.message}`)
      }
    })()
  }

  const handleSortClick = v => {
    onSortChange(v)
  }

  const offerListLoading = loading === 'ask'
  const buyOrderLoading = !buyOrders.data || loading === 'bid'

  return (
    <>
      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="market list table">
          <TableHead className={classes.tableHead}>
            <TableRow>
              <TableHeadCell colSpan={2} padding="none">
                <DashTabs variant="fullWidth" value={tabIndex} onChange={handleTabChange}>
                  <DashTab value={0} label="Offers" badgeContent={offers.total_count} />
                  <DashTab value={1} label="Buy Orders" badgeContent={buyOrders.total_count} />
                </DashTabs>
              </TableHeadCell>
            </TableRow>
          </TableHead>

          {tabIndex === 0 ? (
            <OfferList
              datatable={offers}
              loading={offerListLoading}
              error={error}
              sort={sort}
              onSort={handleSortClick}
              onContact={handleContactClick}
              onRemove={handleRemoveClick}
              currentUserID={currentUserID}
              isMobile={isMobile}
            />
          ) : (
            <OrderList
              datatable={buyOrders}
              loading={buyOrderLoading}
              error={error}
              sort={sort}
              onSort={handleSortClick}
              onContact={handleContactClick}
              onRemove={handleRemoveClick}
              currentUserID={currentUserID}
              isMobile={isMobile}
            />
          )}
        </Table>
      </TableContainer>

      {/* Only display pagination on offer list */}
      {tabIndex === 0 && pagination}

      {tabIndex === 1 && buyOrders.data.length !== 0 && buyOrders.total_count > 10 && (
        <Typography color="textSecondary" align="right" variant="body2" style={{ margin: 8 }}>
          {buyOrders.total_count - 10} more hidden buy orders
          {/* {buyOrders.total_count - 10} more hidden buy orders at &nbsp; */}
          {/* {amount(buyOrders.data[9].price || 0, 'USD')} or less */}
        </Typography>
      )}

      {/* Fixes bottom spacing */}
      {((tabIndex === 0 && offers.total_count === 0) ||
        (tabIndex === 1 && buyOrders.total_count <= 10)) && <div style={{ margin: 8 }}>&nbsp;</div>}

      <ContactDialog
        market={currentMarket}
        open={tabIndex === 0 && !!currentMarket}
        onClose={() => handleContactClick(null)}
      />

      <ContactBuyerDialog
        market={currentMarket}
        open={tabIndex === 1 && !!currentMarket}
        onClose={() => handleContactClick(null)}
      />
    </>
  )
}
MarketList.propTypes = {
  offers: PropTypes.object.isRequired,
  buyOrders: PropTypes.object.isRequired,
  pagination: PropTypes.element,
  error: PropTypes.string,
  loading: PropTypes.bool,
  sort: PropTypes.string,
  tabIndex: PropTypes.number,
  onSortChange: PropTypes.func,
  onTabChange: PropTypes.func,
}
MarketList.defaultProps = {
  pagination: null,
  error: null,
  loading: false,
  sort: null,
  tabIndex: 1,
  onSortChange: () => {},
  onTabChange: () => {},
}

const OfferList = props => {
  const { isMobile } = props
  if (isMobile) {
    return <OfferListMini {...props} />
  }

  return <OfferListDesktop {...props} />
}
OfferList.propTypes = {
  isMobile: PropTypes.bool,
}
OfferList.defaultProps = {
  isMobile: false,
}

const OrderList = props => {
  const { isMobile } = props
  if (isMobile) {
    return <OrderListMini bidMode {...props} />
  }

  return <OrderListDesktop bidMode {...props} />
}
OrderList.propTypes = OfferList.propTypes
OrderList.defaultProps = OfferList.defaultProps

function baseTable(Component) {
  const wrapped = props => {
    const { classes } = useStyles()

    const { currentUserID } = props

    const { onContact, onRemove } = props
    const handleContactClick = marketIdx => {
      onContact(marketIdx)
    }
    const handleRemoveClick = marketIdx => {
      onRemove(marketIdx)
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

    const { datatable, loading, error, bidMode, sort, onSort } = props

    return (
      <>
        <TableBody style={{ opacity: loading ? 0.5 : 1 }}>
          <TableRow>
            <TableHeadCell colSpan={2}>
              <div className={classes.sortButtons}>
                <Chip
                  className={sort === 'price' ? classes.activeSortButtons : null}
                  onClick={() => onSort('price')}
                  label={bidMode ? 'Highest price' : 'Lowest price'}
                  variant="outlined"
                  clickable
                />
                <Chip
                  onClick={() => onSort('best')}
                  className={sort === 'best' ? classes.activeSortButtons : null}
                  label={bidMode ? 'Top buyers' : 'Top sellers'}
                  variant="outlined"
                  clickable
                />
                <Chip
                  onClick={() => onSort('recent')}
                  className={sort === 'recent' ? classes.activeSortButtons : null}
                  label="Recent"
                  variant="outlined"
                  clickable
                />
              </div>
            </TableHeadCell>
          </TableRow>

          {/* <TableRow> */}
          {/*  <TableHeadCell size="small"> */}
          {/*    <Typography color="textSecondary" variant="body2"> */}
          {/*      {bidMode ? 'Buyer' : 'Seller'} */}
          {/*    </Typography> */}
          {/*  </TableHeadCell> */}
          {/*  <TableHeadCell size="small" align="right"> */}
          {/*    <Typography color="textSecondary" variant="body2"> */}
          {/*      {bidMode ? 'Buy Price' : 'Price'} */}
          {/*    </Typography> */}
          {/*  </TableHeadCell> */}
          {/*  {!isMobile && <TableHeadCell size="small" align="center" width={160} />} */}
          {/* </TableRow> */}

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

          {!error && loading && datatable.data.length === 0 && (
            <TableRow>
              <TableCell align="center" colSpan={3}>
                Loading...
              </TableCell>
            </TableRow>
          )}

          {!error && datatable.total_count === 0 && (
            <TableRow>
              <TableCell align="center" colSpan={3}>
                No available {bidMode ? 'orders' : 'offers'}
              </TableCell>
            </TableRow>
          )}

          {datatable.data.map((market, idx) => (
            <TableRow key={market.id} hover>
              <TableCell component="th" scope="row" padding="none">
                <Link href={`/profiles/${market.user.steam_id}`} disableUnderline>
                  <div className={classes.seller}>
                    <Avatar
                      badge={getUserBadgeFromBoons(market.user.boons)}
                      className={classes.avatar}
                      alt={market.user.name}
                      glow={isDonationGlowExpired(market.user.donated_at)}
                      {...retinaSrcSet(market.user.avatar, 40, 40)}
                    />
                    <div>
                      {/* check for redacted data */}
                      {market.user.id ? <strong>{market.user.name}</strong> : <em>████████████</em>}
                      {Boolean(getUserBadgeFromBoons(market.user.boons)) && (
                        <SubscriberBadge
                          style={{ marginLeft: 4 }}
                          type={getUserBadgeFromBoons(market.user.boons)}
                        />
                      )}
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {bidMode ? 'Ordered' : 'Posted'} {dateFromNow(market.created_at)}
                      </Typography>
                      {!bidMode && (
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
                      )}
                    </div>
                  </div>
                </Link>
              </TableCell>
              <Component
                currentUserID={currentUserID}
                market={market}
                onRemove={() => handleRemoveClick(idx)}
                onContact={() => handleContactClick(idx)}
              />
            </TableRow>
          ))}
        </TableBody>

        <VerifiedStatusPopover
          id={popoverElementID}
          open={open}
          anchorEl={anchorEl}
          onClose={handlePopoverClose}
          onMouseEnter={() => debouncePopoverClose.clear()}
          market={datatable.data[currentIndex]}
        />
      </>
    )
  }
  wrapped.propTypes = {
    datatable: PropTypes.object.isRequired,
    error: PropTypes.string,
    loading: PropTypes.bool,
    currentUserID: PropTypes.string,
    isMobile: PropTypes.bool,
    onContact: PropTypes.func,
    onRemove: PropTypes.func,
    bidMode: PropTypes.bool,
  }
  wrapped.defaultProps = {
    error: null,
    loading: false,
    currentUserID: null,
    isMobile: false,
    onContact: () => {},
    onRemove: () => {},
    bidMode: false,
  }

  return wrapped
}

const OfferListDesktop = baseTable(({ market, currentUserID, onRemove, onContact }) => (
  <TableCell align="right">
    <NoSsr>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'right' }}>
        <Typography variant="body2" style={{ marginRight: 16 }}>
          {amount(market.price, market.currency)}
        </Typography>
        {currentUserID === market.user.id ? (
          // HOTFIX! wrapped button on div to prevent mixing up the styles(variant) of 2 buttons.
          <div>
            <Button variant="outlined" onClick={onRemove}>
              Remove
            </Button>
          </div>
        ) : (
          <BuyButton variant="contained" onClick={onContact}>
            Contact Seller
          </BuyButton>
        )}
      </div>
    </NoSsr>
  </TableCell>
))

const OfferListMini = baseTable(({ market, currentUserID, onRemove, onContact }) => (
  <TableCell
    align="right"
    style={{ cursor: 'pointer' }}
    onClick={currentUserID === market.user.id ? onRemove : onContact}>
    <Typography variant="body2">{amount(market.price, market.currency)}</Typography>
    <Typography
      variant="caption"
      color="textSecondary"
      style={{ color: currentUserID === market.user.id ? 'tomato' : '' }}>
      <u>{currentUserID === market.user.id ? 'Remove' : 'View'}</u>
    </Typography>
  </TableCell>
))

const OrderListDesktop = baseTable(({ market, currentUserID, onRemove, onContact }) => (
  <TableCell align="right">
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'right' }}>
      <Typography variant="body2" style={{ marginRight: 16 }}>
        {amount(market.price, market.currency)}
      </Typography>

      {currentUserID === market.user.id ? (
        // HOTFIX! wrapped button on div to prevent mixing up the styles(variant) of 2 buttons.
        <div>
          <Button variant="outlined" onClick={onRemove}>
            Remove
          </Button>
        </div>
      ) : (
        <SellButton
          // Check for redacted user and disable them for opening the dialog.
          disabled={!market.user.id}
          variant="contained"
          onClick={onContact}>
          {market.user.id ? `Contact Buyer` : `Sign in to view`}
        </SellButton>
      )}
    </div>
  </TableCell>
))

const OrderListMini = baseTable(({ market, currentUserID, onRemove, onContact }) => (
  <TableCell
    align="right"
    onClick={() => {
      // Data was redacted, so we can do nothing about it.
      if (!market.user.id) {
        return
      }

      // Logged in user matched th data id, we can invoke remove callback.
      if (currentUserID === market.user.id) {
        onRemove()
        return
      }

      onContact()
    }}
    style={{ cursor: 'pointer' }}>
    <Typography variant="body2" style={{ color: bidColor.A200 }}>
      {amount(market.price, market.currency)}
    </Typography>

    {currentUserID === market.user.id ? (
      <Typography variant="caption" color="textSecondary" style={{ color: 'tomato' }}>
        <u>Remove</u>
      </Typography>
    ) : (
      <Typography variant="caption" color="textSecondary">
        <u>{market.user.id ? 'View' : 'Sign in to view'}</u>
      </Typography>
    )}
  </TableCell>
))
