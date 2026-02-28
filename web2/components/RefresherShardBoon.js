import moment from 'moment'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import { MARKET_BID_EXPR_DAYS } from '@/constants/market'
import Link from '@/components/Link'

export default function RefresherShardBoon({ boons }) {
  if (!boons || boons.indexOf('REFRESHER_SHARD') === -1) {
    return (
      <Box align="center">
        <Typography
          sx={{ color: 'salmon' }}
          component={Link}
          variant="body2"
          href="/expiring-posts">
          This buy order will expires in {MARKET_BID_EXPR_DAYS} days -{' '}
          {moment().add(MARKET_BID_EXPR_DAYS, 'days').calendar()}
        </Typography>
      </Box>
    )
  }

  return (
    <Box align="center">
      <Typography sx={{ color: 'lightgreen' }} component={Link} variant="body2" href="/plus">
        <strong>Refresher Shard</strong>: Automatically refreshes expiring buy orders
      </Typography>
    </Box>
  )
}
