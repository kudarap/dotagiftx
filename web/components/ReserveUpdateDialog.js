import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from 'tss-react/mui'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Typography from '@mui/material/Typography'
import CircularProgress from '@mui/material/CircularProgress'
import TextField from '@mui/material/TextField'
import DeliveredIcon from '@mui/icons-material/AssignmentTurnedIn'
import CancelIcon from '@mui/icons-material/Cancel'
import { STEAM_PROFILE_BASE_URL } from '@/constants/strings'
import { myMarket } from '@/service/api'
import { amount, dateTime } from '@/lib/format'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import {
  MARKET_STATUS_CANCELLED,
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_SOLD,
} from '@/constants/market'
import AppContext from '@/components/AppContext'
import ItemImageDialog from '@/components/ItemImageDialog'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  details: {
    [theme.breakpoints.down('sm')]: {
      display: 'block',
    },
    display: 'inline-flex',
  },
}))

const steamProfileBaseURL = `https://steamcommunity.com'/profiles/`

export default function ReserveUpdateDialog(props) {
  const { classes } = useStyles()
  const { isMobile } = useContext(AppContext)

  const [notes, setNotes] = React.useState('')
  const [error, setError] = React.useState('')
  const [loading, setLoading] = React.useState(false)
  const [loadingCancel, setLoadingCancel] = React.useState(false)

  const { onClose } = props
  const handleClose = () => {
    setNotes('')
    setError('')
    setLoading(false)
    setLoadingCancel(false)
    onClose()
  }

  const { market, onSuccess, onCancel } = props
  const marketUpdate = (payload, setLoader) => {
    if (loading) {
      return
    }

    setLoader(true)
    setError(null)
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, payload)
        if (payload.status === MARKET_STATUS_SOLD) {
          onSuccess()
        } else {
          onCancel()
        }

        handleClose()
      } catch (e) {
        setError(`Error: ${e.message}`)
      }

      setLoader(false)
    })()
  }

  const handleCancelClick = () => {
    marketUpdate(
      {
        status: MARKET_STATUS_CANCELLED,
        notes,
      },
      setLoadingCancel
    )
  }

  const onFormSubmit = evt => {
    evt.preventDefault()

    marketUpdate(
      {
        status: MARKET_STATUS_SOLD,
        notes,
      },
      setLoading
    )
  }

  if (!market) {
    return null
  }

  const { open } = props
  return (
    <Dialog
      fullWidth
      fullScreen={isMobile}
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      <form onSubmit={onFormSubmit}>
        <DialogTitle id="alert-dialog-title">
          Update Reservation
          <DialogCloseButton onClick={handleClose} />
        </DialogTitle>
        <DialogContent>
          <div className={classes.details}>
            <ItemImageDialog item={market.item} />

            <Typography component="h1">
              <Typography component="p" variant="h6">
                {market.item.name}
              </Typography>
              <Typography gutterBottom>
                <Typography color="textSecondary" component="span">
                  {`Status: `}
                </Typography>
                <strong style={{ color: MARKET_STATUS_MAP_COLOR[market.status] }}>
                  {MARKET_STATUS_MAP_TEXT[market.status]}
                </strong>
                <br />
                {market.resell && (
                  <div>
                    <Typography color="textSecondary" component="span">
                      {`Seller Steam ID: `}
                    </Typography>
                    <Link
                      href={steamProfileBaseURL + market.seller_steam_id}
                      target="_blank"
                      rel="noreferrer noopener">
                      {market.seller_steam_id}
                    </Link>
                    <br />
                  </div>
                )}
                <Typography color="textSecondary" component="span">
                  {`Price: `}
                </Typography>
                {amount(market.price, market.currency)}
                <br />
                <Typography color="textSecondary" component="span">
                  {`Reserved: `}
                </Typography>
                {dateTime(market.updated_at)}
                {market.notes && (
                  <>
                    <br />
                    <Typography color="textSecondary" component="span">
                      {`Notes: `}
                    </Typography>
                    <Typography component="ul" variant="body2" style={{ marginTop: 0 }}>
                      {market.notes.split('\n').map(s => (
                        <li>{s}</li>
                      ))}
                    </Typography>
                  </>
                )}
              </Typography>
            </Typography>
          </div>
          <div>
            <TextField
              style={{ marginTop: 16 }}
              InputProps={{ readOnly: true }}
              fullWidth
              color="secondary"
              variant="outlined"
              label="Buyer's Steam profile URL"
              value={`${STEAM_PROFILE_BASE_URL}/${market.partner_steam_id}`}
            />
            <br />
            <br />
            <TextField
              disabled={loading}
              fullWidth
              required
              color="secondary"
              variant="outlined"
              label="Notes"
              helperText="Screenshot URL for verification or reason for cancellation"
              placeholder="https://imgur.com/a/..."
              value={notes}
              onInput={e => setNotes(e.target.value)}
            />
          </div>
        </DialogContent>
        {error && (
          <Typography color="error" align="center" variant="body2">
            {error}
          </Typography>
        )}
        <DialogActions>
          <Button
            disabled={loadingCancel}
            startIcon={loadingCancel ? <CircularProgress size={22} /> : <CancelIcon />}
            onClick={handleCancelClick}
            variant="outlined">
            Cancel Reservation
          </Button>
          <Button
            disabled={loading}
            startIcon={
              loading ? <CircularProgress size={22} color="secondary" /> : <DeliveredIcon />
            }
            variant="outlined"
            color="secondary"
            type="submit">
            {isMobile ? 'Item Delivered' : 'Item Delivered to Buyer'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}
ReserveUpdateDialog.propTypes = {
  market: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
  onCancel: PropTypes.func,
  onSuccess: PropTypes.func,
}
ReserveUpdateDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
  onCancel: () => {},
  onSuccess: () => {},
}
