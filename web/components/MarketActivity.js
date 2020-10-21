import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Divider from '@material-ui/core/Divider'
import { MARKET_STATUS_MAP_COLOR, MARKET_STATUS_MAP_TEXT } from '@/constants/market'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  profile: {
    float: 'left',
    marginRight: theme.spacing(1),
    width: 60,
    height: 60,
  },
  itemImage: { width: 60, height: 40, marginRight: 8, float: 'left' },
}))

export default function MarketActivity({ data, loading }) {
  const classes = useStyles()

  return (
    <>
      {(loading || !data) && <LinearProgress color="secondary" />}
      <ul style={{ paddingLeft: 0, listStyle: 'none', opacity: loading ? 0.5 : 1 }}>
        {data &&
          data.map(market => (
            <li style={{ display: 'flow-root' }}>
              <Divider style={{ margin: '8px 0 8px' }} light />
              <ItemImage
                className={classes.itemImage}
                image={`/200x100/${market.item.image}`}
                title={market.item.name}
                rarity={market.item.rarity}
              />
              <Typography variant="body2">
                <Link href={`/profile/${market.user.steam_id}`} color="textPrimary">
                  {market.user.name}
                </Link>
                &nbsp;
                <span style={{ color: MARKET_STATUS_MAP_COLOR[market.status] }}>
                  {MARKET_STATUS_MAP_TEXT[market.status].toLowerCase()}
                </span>
                &nbsp;
                {market.item.hero}&apos;s&nbsp;
                <Link href={`/${market.item.slug}`} color="secondary">
                  {market.item.name}
                </Link>
                &nbsp;
                {moment(market.updated_at).fromNow()}
              </Typography>
              <Typography
                component="pre"
                color="textSecondary"
                variant="caption"
                style={{ whiteSpace: 'pre-wrap' }}>
                {market.notes}
              </Typography>
            </li>
          ))}
      </ul>
    </>
  )
}
MarketActivity.propTypes = {
  data: PropTypes.arrayOf(PropTypes.object).isRequired,
  loading: PropTypes.bool,
}
MarketActivity.defaultProps = {
  loading: false,
}
