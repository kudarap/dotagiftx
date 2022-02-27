import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import makeStyles from '@mui/styles/makeStyles'
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

const useStyles = makeStyles(theme => ({
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

export default function UserMarketList({ data, loading, error, onSearchInput }) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const [currentMarket, setCurrentMarket] = React.useState(null)
  const handleContactClick = marketIdx => {
    setCurrentMarket(data.data[marketIdx])
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
                    </Link>
                  </TableCell>

                  {!isMobile ? (
                    <>
                      <TableCell align="right">
                        <Typography variant="body2">${market.price.toFixed(2)}</Typography>
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
                      <Typography variant="body2">${market.price.toFixed(2)}</Typography>
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
        market={data.data[currentIndex]}
      />

      <ContactDialog
        market={currentMarket}
        open={!!currentMarket}
        onClose={() => handleContactClick(null)}
      />
    </>
  )
}
UserMarketList.propTypes = {
  onSearchInput: PropTypes.func,
  data: PropTypes.object.isRequired,
  loading: PropTypes.bool,
  error: PropTypes.string,
}
UserMarketList.defaultProps = {
  onSearchInput: () => {},
  loading: false,
  error: null,
}
