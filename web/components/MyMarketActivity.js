import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import lightGreen from '@material-ui/core/colors/lightGreen'
import teal from '@material-ui/core/colors/teal'
import { STEAM_PROFILE_BASE_URL } from '@/constants/strings'
import {
  MARKET_TYPE_ASK,
  MARKET_TYPE_BID,
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_BID_STATUS_MAP_TEXT,
  MARKET_STATUS_LIVE,
  MARKET_STATUS_RESERVED,
  MARKET_STATUS_SOLD,
} from '@/constants/market'
import { VERIFIED_DELIVERY_MAP_ICON, VERIFIED_INVENTORY_MAP_ICON } from '@/constants/verified'
import { amount, daysFromNow } from '@/lib/format'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import { VerifiedStatusPopover } from '@/components/VerifiedStatusCard'

const priceTagStyle = {
  padding: '2px 6px',
  color: 'white',
}

const useStyles = makeStyles(theme => ({
  profile: {
    float: 'left',
    marginRight: theme.spacing(1),
    width: 60,
    height: 60,
  },
  itemImage: { width: 60, height: 40, marginRight: 8, float: 'left' },
  askPriceTag: {
    ...priceTagStyle,
    background: lightGreen[900],
  },
  bidPriceTag: {
    ...priceTagStyle,
    background: teal[900],
  },
  activity: {
    display: 'flow-root',
    borderBottom: `1px ${theme.palette.divider} solid`,
    marginBottom: theme.spacing(1),
    paddingBottom: theme.spacing(1),
  },
  list: {
    padding: theme.spacing(1, 0, 0, 0),
    marginTop: 0,
    borderTop: `1px ${theme.palette.divider} solid`,
    listStyle: 'none',
  },
  text: {
    marginTop: theme.spacing(1),
  },
}))

export default function MyMarketActivity({ datatable, loading, error }) {
  const classes = useStyles()

  const [currentIndex, setIndex] = React.useState(null)
  const [anchorEl, setAnchorEl] = React.useState(null)
  const handlePopoverOpen = event => {
    setIndex(Number(event.currentTarget.dataset.index))
    setAnchorEl(event.currentTarget)
  }
  const handlePopoverClose = () => {
    setAnchorEl(null)
  }
  const open = Boolean(anchorEl)
  const popoverElementID = open ? 'verified-status-popover' : undefined

  if (error) {
    return (
      <Typography className={classes.text} color="error">
        Error {error}
      </Typography>
    )
  }

  if (loading || !datatable.data) {
    return (
      <Typography className={classes.text} color="textSecondary">
        Loading...
      </Typography>
    )
  }

  if (datatable.data.length === 0) {
    return <Typography className={classes.text}>No Result</Typography>
  }

  return (
    <>
      <ul className={classes.list}>
        {datatable.data.map((market, idx) => (
          <li className={classes.activity} key={market.id}>
            <Link href={`/${market.item.slug}`}>
              <ItemImage
                className={classes.itemImage}
                image={market.item.image}
                width={60}
                height={40}
                title={market.item.name}
                rarity={market.item.rarity}
              />
            </Link>
            <Typography variant="body2" color="textSecondary">
              <span style={{ color: MARKET_STATUS_MAP_COLOR[market.status] }}>
                {market.type === MARKET_TYPE_BID
                  ? MARKET_BID_STATUS_MAP_TEXT[market.status]
                  : MARKET_STATUS_MAP_TEXT[market.status]}
              </span>
              <span
                aria-owns={popoverElementID}
                aria-haspopup="true"
                data-index={idx}
                onMouseEnter={handlePopoverOpen}>
                {(market.status === MARKET_STATUS_LIVE ||
                  market.status === MARKET_STATUS_RESERVED) &&
                  VERIFIED_INVENTORY_MAP_ICON[market.inventory_status]}

                {market.status === MARKET_STATUS_SOLD &&
                  VERIFIED_DELIVERY_MAP_ICON[market.delivery_status]}
              </span>
              &nbsp;
              <Link href={`/search?hero=${market.item.hero}`} color="textPrimary">
                {`${market.item.hero}'s`}
              </Link>
              &nbsp;
              <Link href={`/${market.item.slug}`} color="textPrimary">
                {`${market.item.name}`}
              </Link>
              &nbsp;
              {daysFromNow(market.updated_at)}
              &nbsp;for&nbsp;
              <span
                className={
                  market.type === MARKET_TYPE_ASK ? classes.askPriceTag : classes.bidPriceTag
                }>
                {amount(market.price, market.currency)}
              </span>
            </Typography>

            <Typography
              component="pre"
              color="textSecondary"
              variant="caption"
              style={{ whiteSpace: 'pre-wrap', display: 'inline-block' }}>
              {market.partner_steam_id && (
                <Link
                  color="textSecondary"
                  href={`${STEAM_PROFILE_BASE_URL}/${market.partner_steam_id}`}>
                  {`${STEAM_PROFILE_BASE_URL}/${market.partner_steam_id}`}
                  {market.notes && '\n'}
                </Link>
              )}
              {market.notes}
            </Typography>
          </li>
        ))}
      </ul>

      <VerifiedStatusPopover
        id={popoverElementID}
        open={open}
        anchorEl={anchorEl}
        onClose={handlePopoverClose}
        market={datatable.data[currentIndex]}
      />
    </>
  )
}
MyMarketActivity.propTypes = {
  datatable: PropTypes.object.isRequired,
  loading: PropTypes.bool,
  error: PropTypes.string,
}
MyMarketActivity.defaultProps = {
  loading: false,
  error: null,
}
