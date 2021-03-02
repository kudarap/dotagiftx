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
} from '@/constants/market'
import { amount, daysFromNow } from '@/lib/format'
import ItemImage, { retinaSrcSet } from '@/components/ItemImage'
import Link from '@/components/Link'
import Avatar from '@material-ui/core/Avatar'
import AppContext from '@/components/AppContext'

const priceTagStyle = {
  padding: '2px 6px',
  color: 'white',
  // borderRadius: 6,
  // float: 'right',
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
  avatar: {
    float: 'left',
    marginRight: theme.spacing(1),
  },
  activity: {
    display: 'flow-root',
    borderBottom: `1px ${theme.palette.divider} solid`,
    marginBottom: theme.spacing(1),
    paddingBottom: theme.spacing(1),
  },
}))

export default function MyOrderActivity({ datatable, loading, error }) {
  const classes = useStyles()

  const { isMobile } = React.useContext(AppContext)

  if (error) {
    return <p>Error {error}</p>
  }

  if (loading || !datatable.data) {
    return <p>Loading...</p>
  }

  return (
    <>
      <ul style={{ paddingLeft: 0, listStyle: 'none', opacity: loading ? 0.5 : 1 }}>
        {datatable.data.map(market => (
          <li className={classes.activity} key={market.id}>
            {!isMobile && (
              <Avatar
                hidden={isMobile}
                className={classes.avatar}
                {...retinaSrcSet(market.user.avatar, 40, 40)}
                component={Link}
                href={`/profiles/${market.user.steam_id}`}
              />
            )}
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
              You&nbsp;
              <span style={{ color: MARKET_STATUS_MAP_COLOR[market.status] }}>
                {market.type === MARKET_TYPE_BID
                  ? MARKET_BID_STATUS_MAP_TEXT[market.status].toLowerCase()
                  : MARKET_STATUS_MAP_TEXT[market.status].toLowerCase()}
              </span>
              &nbsp;
              <Link href={`/search?hero=${market.item.hero}`} color="textPrimary">
                {`${market.item.hero}'s`}
              </Link>
              &nbsp;
              <Link href={`/${market.item.slug}`} color="textPrimary">
                {`${market.item.name}`}
              </Link>
              &nbsp;from&nbsp;
              <Link href={`/profiles/${market.user.steam_id}`} color="textPrimary">
                {market.user.name}
              </Link>
              &nbsp;
              {daysFromNow(market.updated_at)}
              &nbsp;for&nbsp;
              <span className={classes.bidPriceTag}>{amount(market.price, market.currency)}</span>
            </Typography>

            <Typography
              component="pre"
              color="textSecondary"
              variant="caption"
              style={{ whiteSpace: 'pre-wrap', display: 'inline-block' }}>
              {market.user.steam_id && (
                <Link
                  color="textSecondary"
                  href={`${STEAM_PROFILE_BASE_URL}/${market.user.steam_id}`}>
                  {`${STEAM_PROFILE_BASE_URL}/${market.user.steam_id}`}
                  {market.notes && '\n'}
                </Link>
              )}
              {market.notes}
            </Typography>
          </li>
        ))}
      </ul>
    </>
  )
}
MyOrderActivity.propTypes = {
  datatable: PropTypes.object.isRequired,
  loading: PropTypes.bool,
  error: PropTypes.string,
}
MyOrderActivity.defaultProps = {
  loading: false,
  error: null,
}
