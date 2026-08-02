import React, { useContext } from 'react'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { STEAM_PROFILE_BASE_URL } from '@/constants/strings'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import BidButton from '@/components/BidButton'
import MarketNotes from '@/components/MarketNotes'
import ProfileCard from '@/components/ProfileCard'
import type { Market } from '@/lib/types'

export default function ContactBuyerDialog({
  market,
  open,
  onClose,
}: {
  market: Market | null
  open: boolean
  onClose: () => void
}) {
  const { isMobile } = useContext(AppContext)

  // Check for redacted user and disabled them for opening the dialog.
  if (!market || (market && !market.user.id)) {
    return null
  }

  const storeProfile = `/profiles/${market.user.steam_id}`
  const steamProfileURL = `${STEAM_PROFILE_BASE_URL}/${market.user.steam_id}`

  return (
    <div>
      <Dialog
        fullWidth
        fullScreen={isMobile}
        open={open}
        onClose={onClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
        <DialogTitle id="alert-dialog-title">
          Contact Buyer
          <DialogCloseButton onClick={onClose} />
        </DialogTitle>
        <DialogContent>
          <ProfileCard user={market.user} hideSteamProfile>
            {market.notes && <MarketNotes text={market.notes} />}
          </ProfileCard>

          <Box sx={{ mt: 2 }}>
            <strong>Guides for selling Giftables</strong>
            <Typography
              component="ul"
              variant="body2"
              color="textSecondary"
              style={{ lineHeight: 1.7 }}>
              <li>Please be respectful on the price stated by the buyer.</li>
              <li>Make sure your item exist in your inventory.</li>
              <li>
                Dota 2 Giftables transaction only viable if the two steam user parties have been
                friends for 30 days.
              </li>
              <li>
                Payment agreements will be done between you and the buyer. This website does not
                accept or integrate any payment service.
              </li>
            </Typography>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button component="a" href={storeProfile}>
            Buyer Profile
          </Button>
          <BidButton
            variant="outlined"
            component={Link}
            target="_blank"
            rel="noreferrer noopener"
            disableUnderline
            href={steamProfileURL}>
            Steam Profile
          </BidButton>
        </DialogActions>
      </Dialog>
    </div>
  )
}
