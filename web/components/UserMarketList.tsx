import React, { useContext } from 'react'
import { makeStyles } from 'tss-react/mui'
import { debounce } from '@mui/material'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { lightGreen as green } from '@mui/material/colors'
import {
  VERIFIED_INVENTORY_MAP_ICON,
  VERIFIED_INVENTORY_VERIFIED_RESELL,
} from '@/constants/verified'
import Link from '@/components/Link'
import BuyButton from '@/components/BuyButton'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'
import ContactDialog from '@/components/ContactDialog'
import TableSearchInput from '@/components/TableSearchInput'
import AppContext from '@/components/AppContext'
import { VerifiedStatusPopover } from '@/components/VerifiedStatusCard'
import type { ActivityMarket } from '@/components/MarketActivity'
import type { MarketDatatable } from '@/components/MarketList'

const useStyles = makeStyles()(theme => ({
  seller: {
    display: 'inline-flex',
  },
  link: {
    padding: theme.spacing(2, 2, 2, 0),
    display: 'flex',
  },
  image: {
    margin: theme.spacing(-1, 1, -1, 1),
    width: 77,
    height: 55,
  },
  buyText: {
    color: green[600],
  },
}))

interface UserMarketListProps {
  data: MarketDatatable
  loading?: boolean
  error?: string | null
  onSearchInput?: (value: string) => void
  searchValue?: string
}

export default function UserMarketList({
  data,
  loading,
  error,
  onSearchInput = () => {},
  searchValue = '',
}: UserMarketListProps) {
  const { classes } = useStyles()
  const { isMobile } = useContext(AppContext)

  const [currentMarket, setCurrentMarket] = React.useState<ActivityMarket | null>(null)
  const handleContactClick = (marketIdx: number | null) => {
    setCurrentMarket(marketIdx === null ? null : data.data[marketIdx])
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

  return (
    <>
      <TableContainer component={Paper}>
        <Table aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableHeadCell padding="none" colSpan={isMobile ? 2 : 1}>
                <TableSearchInput
                  fullWidth
                  loading={loading}
                  onInput={onSearchInput}
                  value={searchValue}
                  color="secondary"
                  placeholder="Filter user items"
                />
              </TableHeadCell>
              {!isMobile && (
                <>
                  <TableHeadCell align="right">Price</TableHeadCell>
                  <TableHeadCell align="right" width={156} />
                </>
              )}
            </TableRow>
          </TableHead>
          <TableBody style={loading ? { opacity: 0.5 } : undefined}>
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

            {!error && !data && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  Loading...
                </TableCell>
              </TableRow>
            )}

            {data.data &&
              data.data.map((market, idx) => (
                <TableRow key={market.id} hover>
                  <TableCell component="th" scope="row" padding="none">
                    <Link className={classes.link} href={`/${market.item.slug}`} disableUnderline>
                      <ItemImage
                        className={classes.image}
                        image={market.item.image || ''}
                        width={77}
                        height={55}
                        title={market.item.name}
                        rarity={market.item.rarity || 'regular'}
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
                            : VERIFIED_INVENTORY_MAP_ICON[market.inventory_status!]}
                        </span>
                        <br />
                        <Typography variant="caption" color="textSecondary">
                          {market.item.hero}
                        </Typography>
                        <RarityTag rarity={market.item.rarity || 'regular'} />
                      </div>
                    </Link>
                  </TableCell>

                  {!isMobile ? (
                    <>
                      <TableCell align="right">
                        <Typography variant="body2">
                          ${(market.price || 0).toFixed(2)}
                        </Typography>
                      </TableCell>
                      <TableCell align="right">
                        <BuyButton variant="contained" onClick={() => handleContactClick(idx)}>
                          Contact Seller
                        </BuyButton>
                      </TableCell>
                    </>
                  ) : (
                    <TableCell
                      align="right"
                      onClick={() => handleContactClick(idx)}
                      style={{ cursor: 'pointer' }}>
                      <Typography variant="body2">
                        ${(market.price || 0).toFixed(2)}
                      </Typography>
                      <Typography variant="caption" color="textSecondary">
                        <u>View</u>
                      </Typography>
                    </TableCell>
                  )}
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>

      <VerifiedStatusPopover
        id={popoverElementID}
        open={open}
        anchorEl={anchorEl}
        onClose={handlePopoverClose}
        onMouseEnter={() => debouncePopoverClose.clear()}
        market={currentIndex === null ? null : data.data[currentIndex]}
      />

      <ContactDialog
        market={currentMarket}
        open={!!currentMarket}
        onClose={() => handleContactClick(null)}
      />
    </>
  )
}
