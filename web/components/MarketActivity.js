import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Divider from '@material-ui/core/Divider'
import { Link } from '@material-ui/core'
import { MARKET_STATUS_MAP_COLOR, MARKET_STATUS_MAP_TEXT } from '@/constants/market'
import ItemImage from '@/components/ItemImage'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  profile: {
    float: 'left',
    marginRight: theme.spacing(1),
    width: 60,
    height: 60,
  },
  itemImage: { width: 60, height: 40, marginRight: 8, float: 'left' },
}))

export default function MarketActivity({ data }) {
  const classes = useStyles()

  return (
    <ul style={{ paddingLeft: 0, listStyle: 'none' }}>
      {data.map(market => (
        <li>
          <Divider style={{ margin: '8px 0 8px' }} light />
          <ItemImage
            className={classes.itemImage}
            image={`/200x100/${market.item.image}`}
            title={market.item.name}
            rarity={market.item.rarity}
          />
          <Typography variant="body2">
            {market.user.name}{' '}
            <span style={{ color: MARKET_STATUS_MAP_COLOR[market.status] }}>
              {MARKET_STATUS_MAP_TEXT[market.status].toLowerCase()}
            </span>
            &nbsp;
            {market.item.hero}&apos;s&nbsp;
            <Link href={`/item/${market.item.slug}`} color="secondary">
              {market.item.name}
            </Link>
            &nbsp;
            {moment(market.updated_at).fromNow()}
          </Typography>
          <Typography component="pre" color="textSecondary" variant="caption">
            {market.notes}
          </Typography>
        </li>
      ))}
    </ul>
  )
}
MarketActivity.propTypes = {
  data: PropTypes.arrayOf(PropTypes.object).isRequired,
}
