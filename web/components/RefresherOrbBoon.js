import { addDays, format } from 'date-fns'
import PropTypes from 'prop-types'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import { MARKET_ASK_EXPR_DAYS } from '@/constants/market'
import Link from '@/components/Link'

export default function RefresherOrbBoon({ boons }) {
  if (!boons || boons.indexOf('REFRESHER_ORB') === -1) {
    return (
      <Box align="center">
        <Typography
          sx={{ color: 'salmon' }}
          component={Link}
          variant="body2"
          href="/expiring-posts"
        >
          This listing will expires in {MARKET_ASK_EXPR_DAYS} days -{' '}
          {format(addDays(new Date(), MARKET_ASK_EXPR_DAYS, 'days'), 'MMM d')}
        </Typography>
      </Box>
    )
  }

  return (
    <Box align="center">
      <Typography sx={{ color: 'lightgreen' }} component={Link} variant="body2" href="/plus">
        <strong>Refresher Orb</strong>: Automatically refreshes expiring buy orders and listings
      </Typography>
    </Box>
  )
}

RefresherOrbBoon.propTypes = {
  boons: PropTypes.string,
}

RefresherOrbBoon.defaultProps = {
  boons: '',
}
