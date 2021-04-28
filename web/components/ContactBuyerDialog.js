import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import Avatar from '@/components/Avatar'
import ChipLink from '@/components/ChipLink'
import {
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import { USER_STATUS_MAP_TEXT } from '@/constants/user'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import BidButton from '@/components/BidButton'
import MarketNotes from '@/components/MarketNotes'
import DonatorBadge from '@/components/DonatorBadge'

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  profileName: {
    [theme.breakpoints.down('xs')]: {
      fontSize: theme.typography.h6.fontSize,
    },
  },
  avatar: {
    [theme.breakpoints.down('xs')]: {
      margin: '0 auto',
    },
    width: 100,
    height: 100,
    marginRight: theme.spacing(1.5),
  },
}))

export default function ContactBuyerDialog(props) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const { market, open, onClose } = props

  // Check for redacted user and disabled them for opening the dialog.
  if (!market || (market && !market.user.id)) {
    return null
  }

  const storeProfile = `/profiles/${market.user.steam_id}`
  const steamProfileURL = `${STEAM_PROFILE_BASE_URL}/${market.user.steam_id}`
  const dota2Inventory = `${steamProfileURL}/inventory#570`

  const isProfileReported = Boolean(market.user.status)

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
          <div
            className={classes.details}
            style={
              isProfileReported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : null
            }>
            <a href={storeProfile} target="_blank" rel="noreferrer noopener">
              <Avatar
                className={classes.avatar}
                glow={Boolean(market.user.donation)}
                {...retinaSrcSet(market.user.avatar, 100, 100)}
              />
            </a>
            <Typography component="h1">
              <Typography
                className={classes.profileName}
                component="p"
                variant="h4"
                color={isProfileReported ? 'error' : ''}>
                {market.user.name}
                {Boolean(market.user.donation) && (
                  <DonatorBadge
                    style={{ marginLeft: 4, marginTop: 12, position: 'absolute' }}
                    size="medium">
                    DONATOR
                  </DonatorBadge>
                )}
              </Typography>
              {isProfileReported && (
                <Typography color="error">{USER_STATUS_MAP_TEXT[market.user.status]}</Typography>
              )}
              <Typography gutterBottom>
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
              </Typography>
            </Typography>
          </div>

          <Typography variant="body2" color="textSecondary">
            <br />
            Guides for selling Giftables
            <ul>
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
            </ul>
          </Typography>
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
            Check Steam Profile
          </BidButton>
        </DialogActions>
      </Dialog>
    </div>
  )
}
ContactBuyerDialog.propTypes = {
  market: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
ContactBuyerDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
}
