import { MARKET_BID_EXPR_DAYS } from '@/constants/market'
import moment from 'moment'
import Link from '@/components/Link'
import { Typography } from '@mui/material'

export default function RefresherShardBoon({ boons }) {
  if (!boons || boons.indexOf('REFRESHER_SHARD') === -1) {
    return (
      <div align="center">
        <Typography
          sx={{ color: 'salmon' }}
          component={Link}
          variant="body2"
          href="/expiring-posts">
          This buy order will expires in {MARKET_BID_EXPR_DAYS} days -{' '}
          {moment().add(MARKET_BID_EXPR_DAYS, 'days').calendar()}
        </Typography>
      </div>
    )
  }

  return (
    <div align="center">
      <Typography sx={{ color: 'lightgreen' }} component={Link} variant="body2" href="/plus">
        <strong>Refresher Shard</strong>: Automatically refreshes expiring buy orders
      </Typography>
    </div>
  )
}
