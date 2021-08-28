import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'

import FormGroup from '@material-ui/core/FormGroup'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import Checkbox from '@material-ui/core/Checkbox'
import CheckBoxOutlineBlankIcon from '@material-ui/icons/CheckBoxOutlineBlank'
import CheckBoxIcon from '@material-ui/icons/CheckBox'
import Favorite from '@material-ui/icons/Favorite'
import FavoriteBorder from '@material-ui/icons/FavoriteBorder'

import { amount, dateCalendar } from '@/lib/format'
import { STEAM_PROFILE_BASE_URL } from '@/constants/strings'
import DialogCloseButton from '@/components/DialogCloseButton'
import {
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_TYPE_BID,
} from '@/constants/market'
import AppContext from '@/components/AppContext'
import { TextField } from '@material-ui/core'
import ItemImageDialog from '@/components/ItemImageDialog'

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      display: 'block',
    },
    display: 'inline-flex',
  },
}))

export default function BuyingGuildeDialog(props) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const { open, onClose } = props

  const handleClose = () => {
    onClose()
  }

  const [state, setState] = React.useState({
    checkedA: false,
    checkedB: false,
    checkedC: false,
    checkedD: false,
    checkedF: false,
    checkedG: false,
  })

  const handleChange = event => {
    setState({ ...state, [event.target.name]: event.target.checked })
  }

  return (
    <Dialog
      fullWidth
      fullScreen={isMobile}
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      <DialogTitle id="alert-dialog-title">
        Buying Safety Tips
        <DialogCloseButton onClick={handleClose} />
      </DialogTitle>
      <DialogContent>
        <Typography component="h1">
          Please check all the boxes and make sure you agree with it
        </Typography>

        <FormGroup>
          <FormControlLabel
            control={<Checkbox checked={state.checkedA} onChange={handleChange} name="checkedA" />}
            label="Buyer profile details are only visible to signed in website users."
          />
          <FormControlLabel
            control={<Checkbox checked={state.checkedB} onChange={handleChange} name="checkedB" />}
            label="As Giftables involves a party having to go first, please always check seller's reputation through SteamRep."
          />
          <FormControlLabel
            control={<Checkbox checked={state.checkedC} onChange={handleChange} name="checkedC" />}
            label="Payment agreements will be done between you and the seller. This website does not accept or integrate any payment service."
          />
          <FormControlLabel
            control={<Checkbox checked={state.checkedD} onChange={handleChange} name="checkedD" />}
            label="Please kindly remove this buy order after use to prevent seller's contacting you."
          />
          <FormControlLabel
            control={<Checkbox checked={state.checkedF} onChange={handleChange} name="checkedF" />}
            label="Some guide when got scammed should put here link"
          />
        </FormGroup>
      </DialogContent>
    </Dialog>
  )
}
BuyingGuildeDialog.propTypes = {
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
BuyingGuildeDialog.defaultProps = {
  open: false,
  onClose: () => {},
}
