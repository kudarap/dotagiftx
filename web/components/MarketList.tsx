import React, { useContext } from 'react'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import { Box, debounce, NoSsr, Tooltip } from '@mui/material'
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
import moment from 'moment'
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
import SubscriberBadge from '@/components/SubscriberBadge'
import type { SubscriberBadgeType } from '@/components/SubscriberBadge'
import { getUserBadgeFromBoons } from '@/lib/badge'
import type { ActivityMarket } from '@/components/MarketActivity'

const displayProfileJoinedDate = false
const displayPostId = false

const useStyles = makeStyles()(theme => ({
  seller: {
    display: 'flex',
    padding: theme.spacing(2),
  },
  avatar: {
    marginRight: theme.spacing(1),
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

export interface MarketDatatable {
  data: ActivityMarket[]
  total_count: number
}

interface MarketListProps {
  offers: MarketDatatable
  buyOrders: MarketDatatable
  pagination?: React.ReactNode
  error?: string | null
  loadingType?: string | null
  sort?: string | null
  tabIndex?: number
  onSortChange?: (v: string) => void
  onTabChange?: (v: number) => void
}

export default function MarketList({
  offers,
  buyOrders,
  error,
  loadingType,
  sort,
  pagination,
  tabIndex = 1,
  onSortChange = () => {},
  onTabChange = () => {},
}: MarketListProps) {
  const { classes } = useStyles()
  const { isMobile, currentAuth } = useContext(AppContext)
  const currentUserID = currentAuth?.user_id || null

  const router = useRouter()
  const handleTabChange = (e: React.SyntheticEvent, value: number) => {
    onTabChange(value)
  }

  const [currentMarket, setCurrentMarket] = React.useState<ActivityMarket | null>(null)
  const handleContactClick = (marketIdx: number | null) => {
    let src = offers
    if (tabIndex === 1) {
      src = buyOrders
    }

    setCurrentMarket(marketIdx === null ? null : src.data[marketIdx])
  }
  const handleRemoveClick = (marketIdx: number) => {
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
        console.error(`Error: ${(e as Error).message}`)
      }
    })()
  }
  const handleSortClick = (v: string) => {
    onSortChange(v)
  }

  const offerListLoading = loadingType === 'ask'
  const buyOrderLoading = !buyOrders.data || loadingType === 'bid'

  return (
    <>
      <TableContainer component={Paper}>
        <Table aria-label="market list table">
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

function OfferList(props: BaseTableProps) {
  const { isMobile } = props
  if (isMobile) {
    return <OfferListMini {...props} />
  }

  return <OfferListDesktop {...props} />
}

function OrderList(props: BaseTableProps) {
  const { isMobile } = props
  if (isMobile) {
    return <OrderListMini bidMode {...props} />
  }

  return <OrderListDesktop bidMode {...props} />
}

interface BaseTableProps {
  bidMode?: boolean
  currentUserID: string | number | null
  datatable: MarketDatatable
  error?: string | null
  isMobile?: boolean
  loading?: boolean
  onContact: (idx: number | null) => void
  onRemove: (idx: number) => void
  sort?: string | null
  onSort?: (v: string) => void
}

type ListActionProps = Omit<
  BaseTableProps,
  | 'datatable'
  | 'error'
  | 'isMobile'
  | 'loading'
  | 'sort'
  | 'onSort'
  | 'bidMode'
  | 'onRemove'
  | 'onContact'
> & {
  market: ActivityMarket
  onRemove: () => void
  onContact: () => void
}

function baseTable(Component: React.ComponentType<ListActionProps>) {
  function Wrapped(props: BaseTableProps) {
    const { classes } = useStyles()

    const { currentUserID } = props

    const { onContact, onRemove } = props
    const handleContactClick = (marketIdx: number) => {
      onContact(marketIdx)
    }
    const handleRemoveClick = (marketIdx: number) => {
      onRemove(marketIdx)
    }

    const [currentIndex, setIndex] = React.useState<number | null>(null)
    const [anchorEl, setAnchorEl] = React.useState<HTMLElement | null>(null)
    const debouncePopoverClose = debounce(() => {
      setAnchorEl(null)
      setIndex(null)
    }, 150)
    const handlePopoverOpen = (event: React.MouseEvent<HTMLElement>) => {
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
                  className={sort === 'price' ? classes.activeSortButtons : undefined}
                  onClick={() => onSort?.('price')}
                  label={bidMode ? 'Highest price' : 'Lowest price'}
                  variant="outlined"
                  clickable
                />
                <Chip
                  onClick={() => onSort?.('recent')}
                  className={sort === 'recent' ? classes.activeSortButtons : undefined}
                  label="Recent"
                  variant="outlined"
                  clickable
                />
                <Chip
                  onClick={() => onSort?.('best')}
                  className={sort === 'best' ? classes.activeSortButtons : undefined}
                  label={bidMode ? 'Top buyers' : 'Top sellers'}
                  variant="outlined"
                  clickable
                />
              </div>
            </TableHeadCell>
          </TableRow>

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
                <div className={classes.seller}>
                  <Link
                    disableUnderline
                    disabled={!market.user.id}
                    href={`/profiles/${market.user.steam_id}`}
                  >
                    <Avatar
                      badge={
                        (getUserBadgeFromBoons(market.user.boons) ||
                          undefined) as SubscriberBadgeType
                      }
                      className={classes.avatar}
                      alt={market.user.name}
                      glow={isDonationGlowExpired(market.user.donated_at || '')}
                      {...retinaSrcSet(market.user.avatar, 40, 40)}
                    />
                  </Link>
                  <div>
                    <strong>{market.user.name}</strong>
                    {market.user.id && Boolean(getUserBadgeFromBoons(market.user.boons)) && (
                      <SubscriberBadge
                        style={{ marginLeft: '0.375rem' }}
                        type={
                          (getUserBadgeFromBoons(market.user.boons) ||
                            undefined) as SubscriberBadgeType
                        }
                      />
                    )}
                    {displayProfileJoinedDate && (
                      <Typography variant="caption" sx={{ ml: 0.5 }}>
                        joined {moment(market.user.created_at).fromNow()}
                      </Typography>
                    )}
                    <br />

                    {/* Copy id to clipboard */}
                    {displayPostId && (
                      <>
                        <Tooltip placement="bottom" title="Copy id to clipboard" arrow>
                          <Typography
                            variant="caption"
                            color="textSecondary"
                            style={{ zIndex: 100 }}
                          >
                            {market.id.split('-')[0]}
                          </Typography>
                        </Tooltip>
                        &nbsp;&middot;&nbsp;
                      </>
                    )}

                    <Typography variant="caption" color="textSecondary">
                      {bidMode ? 'Ordered' : 'Posted'} {dateFromNow(market.created_at || '')}
                    </Typography>

                    {/* Verification status badge with Tooltip */}
                    {!bidMode && (
                      <Box
                        component="span"
                        sx={{ cursor: 'pointer' }}
                        data-index={idx}
                        aria-owns={popoverElementID}
                        aria-haspopup="true"
                        onMouseLeave={debouncePopoverClose}
                        onMouseEnter={handlePopoverOpen}
                      >
                        {market.resell
                          ? VERIFIED_INVENTORY_MAP_ICON[VERIFIED_INVENTORY_VERIFIED_RESELL]
                          : VERIFIED_INVENTORY_MAP_ICON[market.inventory_status!]}
                      </Box>
                    )}
                  </div>
                </div>
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
          market={currentIndex === null ? null : datatable.data[currentIndex]}
          anchorEl={anchorEl}
          onClose={handlePopoverClose}
          onMouseEnter={() => debouncePopoverClose.clear()}
        />
      </>
    )
  }

  return Wrapped
}

const OfferListDesktop = baseTable(
  ({ market, currentUserID, onRemove, onContact }: ListActionProps) => (
    <TableCell align="right">
      <NoSsr>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'right' }}>
          <Typography variant="body2" style={{ marginRight: 16 }}>
            {amount(market.price || 0, market.currency)}
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
  )
)

const OfferListMini = baseTable(
  ({ market, currentUserID, onRemove, onContact }: ListActionProps) => (
    <TableCell
      align="right"
      style={{ cursor: 'pointer' }}
      onClick={currentUserID === market.user.id ? onRemove : onContact}
    >
      <Typography variant="body2">{amount(market.price || 0, market.currency)}</Typography>
      <Typography
        variant="caption"
        color="textSecondary"
        style={{ color: currentUserID === market.user.id ? 'tomato' : '' }}
      >
        <u>{currentUserID === market.user.id ? 'Remove' : 'View'}</u>
      </Typography>
    </TableCell>
  )
)

const OrderListDesktop = baseTable(
  ({ market, currentUserID, onRemove, onContact }: ListActionProps) => (
    <TableCell align="right">
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'right' }}>
        <Typography variant="body2" style={{ marginRight: 16 }}>
          {amount(market.price || 0, market.currency)}
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
            onClick={onContact}
          >
            {market.user.id ? `Contact Buyer` : `Sign in to view`}
          </SellButton>
        )}
      </div>
    </TableCell>
  )
)

const OrderListMini = baseTable(
  ({ market, currentUserID, onRemove, onContact }: ListActionProps) => (
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
      style={{ cursor: 'pointer' }}
    >
      <Typography variant="body2" style={{ color: bidColor.A200 }}>
        {amount(market.price || 0, market.currency)}
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
  )
)
