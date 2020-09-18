import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import { Avatar } from '@material-ui/core'
import { CDN_URL } from '@/service/api'
import ChipLink from '@/components/ChipLink'
import { STEAM_PROFILE_BASE_URL, STEAMREP_PROFILE_BASE_URL } from '@/constants/strings'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
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

export default function ContactDialog(props) {
  const classes = useStyles()

  const { market, open, onClose } = props

  if (!market) {
    return null
  }

  const storeProfile = `/user/${market.user.steam_id}`
  const steamProfileURL = `${STEAM_PROFILE_BASE_URL}/${market.user.steam_id}`
  const dota2Inventory = `${steamProfileURL}/inventory#570`

  return (
    <div>
      <Dialog
        fullWidth
        open={open}
        onClose={onClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
        <DialogTitle id="alert-dialog-title">
          Contact Seller
          <DialogCloseButton onClick={onClose} />
        </DialogTitle>
        <DialogContent>
          <div className={classes.details}>
            <a href={storeProfile} target="_blank" rel="noreferrer noopener">
              <Avatar className={classes.avatar} src={`${CDN_URL}/${market.user.avatar}`} />
            </a>
            <Typography component="h1">
              <Typography component="p" variant="h4">
                {market.user.name}
              </Typography>
              <Typography gutterBottom>
                <Typography color="textSecondary" component="span">
                  {`quick links: `}
                </Typography>
                {/* <ChipLink label="Steam Profile" href={steamProfileURL} /> */}
                {/* &nbsp; */}
                <ChipLink
                  label="SteamRep"
                  href={`${STEAMREP_PROFILE_BASE_URL}/${market.user.steam_id}`}
                />
                &nbsp;
                <ChipLink label="Steam Inventory" href={dota2Inventory} />
                {market.notes && (
                  <>
                    <br />
                    <Typography color="textSecondary" component="span">
                      {`notes: `}
                    </Typography>
                    {market.notes}
                  </>
                )}
              </Typography>
            </Typography>
          </div>

          <div>
            <br />
            Guides for Giftables
            <ul>
              <li>
                Dota 2 giftables transaction only viable if the two steam user parties have been
                friends for 30 days.
              </li>
              <li>
                As giftables involves a party having to go first, please always check seller&apos;s
                reputation through&nbsp;
                <Link
                  disableUnderline={false}
                  href={`${STEAMREP_PROFILE_BASE_URL}/${market.user.steam_id}`}
                  target="_blank"
                  rel="noreferrer noopener">
                  SteamRep
                </Link>
                .
              </li>
              <li>
                Always check the item/set availability on seller&apos;s Dota 2 {` `}
                <Link
                  disableUnderline={false}
                  href={dota2Inventory}
                  target="_blank"
                  rel="noreferrer noopener">
                  inventory
                </Link>
                .
              </li>
            </ul>
          </div>
        </DialogContent>
        <DialogActions>
          <Button component="a" href={storeProfile}>
            View Seller Items
          </Button>
          <Button
            color="secondary"
            variant="outlined"
            component={Link}
            target="_blank"
            rel="noreferrer noopener"
            disableUnderline
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
