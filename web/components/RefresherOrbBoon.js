import { MARKET_ASK_EXPR_DAYS } from '@/constants/market'
import moment from 'moment'
import Link from '@/components/Link'
import { Typography } from '@mui/material'

export default function RefresherOrbBoon({ boons }) {
  if (!boons || boons.indexOf('REFRESHER_ORB') === -1) {
    return (
      <div align="center">
        <Typography
          sx={{ color: 'salmon' }}
          component={Link}
          variant="body2"
          href="/expiring-posts">
          This listing will expires in {MARKET_ASK_EXPR_DAYS} days -{' '}
          {moment().add(MARKET_ASK_EXPR_DAYS, 'days').calendar()}
        </Typography>
      </div>
    )
  }

  return (
    <div align="center">
      <Typography sx={{ color: 'lightgreen' }} component={Link} variant="body2" href="/plus">
        <strong>Refresher Orb</strong>: Automatically refreshes expiring buy orders and listings
      </Typography>
    </div>
  )
}
