import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Typography from '@mui/material/Typography'
import ChipLink from '@/components/ChipLink'
import {
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import MarketNotes from '@/components/MarketNotes'
import ProfileCard from '@/components/ProfileCard'

export default function ContactDialog(props) {
  const { isMobile } = useContext(AppContext)

  const { market, open, onClose } = props
  if (!market) {
    return null
  }

  const storeProfile = `/profiles/${market.user.steam_id}`
  const steamProfileURL = `${STEAM_PROFILE_BASE_URL}/${market.user.steam_id}`
  const dota2Inventory = `${steamProfileURL}/inventory#570`

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
          Contact Seller
          <DialogCloseButton onClick={onClose} />
        </DialogTitle>
        <DialogContent>
          <ProfileCard user={market.user}>
            {/* <Typography color="textSecondary" component="span"> */}
            {/*  {`Links: `} */}
            {/* </Typography> */}
            {/* <ChipLink label="Steam Profile" href={steamProfileURL} /> */}
            {/* &nbsp; */}
            <ChipLink
              label="SteamRep"
              href={`${STEAMREP_PROFILE_BASE_URL}/${market.user.steam_id}`}
            />
            &nbsp;
            <ChipLink
              label="Dotabuff"
              href={`${DOTABUFF_PROFILE_BASE_URL}/${market.user.steam_id}`}
            />
            &nbsp;
            <ChipLink label="Steam Inventory" href={dota2Inventory} />
            {market.notes && <MarketNotes text={market.notes} />}
          </ProfileCard>

          <Typography variant="body2" color="textSecondary">
            <br />
            Guides for buying Giftables
            <ul style={{ lineHeight: 1.7 }}>
              <li>
                Always check the item or set availability on seller&apos;s Dota 2 {` `}
                <Link
                  style={{ textDecoration: 'underline' }}
                  href={dota2Inventory}
                  target="_blank"
                  rel="noreferrer noopener">
                  inventory
                </Link>
                .
              </li>
              <li>
                Ask seller to{' '}
                <Typography variant="inherit" component="span" color="white">
                  Reserve
                </Typography>{' '}
                the item to your profile. This can be use later for{' '}
                <Typography variant="inherit" component="span" color="white">
                  Scam Report
                </Typography>{' '}
                and{' '}
                <Link style={{ textDecoration: 'underline' }} href="/my-orders">
                  order details
                </Link>{' '}
                to avoid{' '}
                <Typography variant="inherit" component="span" color="white">
                  impersonation
                </Typography>
                .<sup style={{ color: 'yellowgreen' }}>NEW</sup>
              </li>
              <li>
                Dota 2 Giftables transaction only viable if the two steam user parties have been
                friends for 30 days.
              </li>
              <li>
                As Giftables involves a party having to go first, please always check seller&apos;s
                reputation through&nbsp;
                <Link
                  style={{ textDecoration: 'underline' }}
                  href={`${STEAMREP_PROFILE_BASE_URL}/${market.user.steam_id}`}
                  target="_blank"
                  rel="noreferrer noopener">
                  SteamRep
                </Link>
                &nbsp;and{' '}
                <Link
                  style={{ textDecoration: 'underline' }}
                  href={`/profiles/${market.user.steam_id}/delivered`}>
                  transaction history
                </Link>
                .
              </li>

              <li>
                If you need a middleman, I only suggest you get{' '}
                <Link href="/middlemen" target="_blank" color="secondary">
                  Middleman here
                </Link>
                .
              </li>

              {/* <li> */}
              {/*  Official SteamRep middleman may assist in middle manning for the trade, or{' '} */}
              {/*  <Link */}
              {/*    style={{ textDecoration: 'underline' }} */}
              {/*    href="https://www.reddit.com/r/dota2trade/" */}
              {/*    target="_blank" */}
              {/*    rel="noreferrer noopener"> */}
              {/*    r/Dota2Trade */}
              {/*  </Link>{' '} */}
              {/*  mod may assist as well in this. */}
              {/* </li> */}
            </ul>
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button component="a" href={storeProfile}>
            View Seller Items
          </Button>
          <Button
            color="secondary"
            variant="outlined"
            component={Link}
            disableUnderline
            target="_blank"
            rel="noreferrer noopener"
            href={steamProfileURL}>
            Check Steam Profile
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  )
}
ContactDialog.propTypes = {
  market: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
ContactDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
}
