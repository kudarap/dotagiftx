import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import startsWith from 'lodash/startsWith'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import TextField from '@material-ui/core/TextField'
import CircularProgress from '@material-ui/core/CircularProgress'
import ReserveIcon from '@material-ui/icons/EventAvailable'
import RemoveIcon from '@material-ui/icons/Delete'
import * as url from '@/lib/url'
import { amount, dateTime } from '@/lib/format'
import { myMarket } from '@/service/api'
import Button from '@/components/Button'
import Link from '@/components/Link'
import DialogCloseButton from '@/components/DialogCloseButton'
import {
  MARKET_STATUS_BID_COMPLETED,
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_REMOVED,
} from '@/constants/market'
import AppContext from '@/components/AppContext'
import ItemImageDialog from '@/components/ItemImageDialog'

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      display: 'block',
    },
    display: 'inline-flex',
  },
}))

const steamCommunityBaseURL = 'https://steamcommunity.com'

export default function BuyOrderUpdateDialog(props) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const [steamProfileURL, setSteamProfileURL] = React.useState('')
  const [notes, setNotes] = React.useState('')
  const [error, setError] = React.useState('')
  const [loading, setLoading] = React.useState(false)

  const { onClose } = props
  const handleClose = () => {
    setSteamProfileURL('')
    setNotes('')
    setError('')
    setLoading(false)
    onClose()
  }

  const { onRemove } = props
  const handleRemove = () => {
    onRemove()
    handleClose()
  }

  const { market } = props
  const handleRemoveClick = () => {
    setLoading(true)
    setError(null)
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, { status: MARKET_STATUS_REMOVED })
        handleRemove()
      } catch (e) {
        setError(`Error: ${e.message}`)
      }

      setLoading(false)
    })()
  }

  const { onSuccess } = props
  const onFormSubmit = evt => {
    evt.preventDefault()

    if (!url.isValid(steamProfileURL)) {
      setError('Steam Profile is not a valid URL.')
      return
    }
    if (!startsWith(steamProfileURL, steamCommunityBaseURL, 0)) {
      setError(`Steam Profile should start with ${steamCommunityBaseURL}`)
      return
    }

    const payload = {
      status: MARKET_STATUS_BID_COMPLETED,
      partner_steam_id: steamProfileURL,
      notes,
    }

    setLoading(true)
    setError(null)
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, payload)
        handleClose()
        onSuccess()
      } catch (e) {
        setError(`Error: ${e.message}`)
      }

      setLoading(false)
    })()
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
          Update Buy Order
          <DialogCloseButton onClick={handleClose} />
        </DialogTitle>
        <DialogContent>
          <div className={classes.details}>
            <ItemImageDialog item={market.item} />

            <Typography component="h1">
              <Typography variant="h6" component={Link} href="/[slug]" as={`/${market.item.slug}`}>
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
                <Typography color="textSecondary" component="span">
                  {`Price: `}
                </Typography>
                {amount(market.price, market.currency)}
                <br />
                <Typography color="textSecondary" component="span">
                  {`Listed: `}
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
              style={{ marginTop: 8 }}
              disabled={loading}
              fullWidth
              required
              color="secondary"
              variant="outlined"
              label="Seller's Steam profile URL"
              helperText="Records seller history for tracking good and bad reputation."
              placeholder="https://steamcommunity.com/..."
              value={steamProfileURL}
              onInput={e => setSteamProfileURL(e.target.value)}
            />
            <br />
            <br />
            <TextField
              disabled={loading}
              fullWidth
              color="secondary"
              variant="outlined"
              label="Notes"
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
            disabled={loading}
            startIcon={<RemoveIcon />}
            onClick={handleRemoveClick}
            variant="outlined">
            {isMobile ? 'Remove' : 'Remove order'}
          </Button>
          <Button
            startIcon={loading ? <CircularProgress size={22} color="secondary" /> : <ReserveIcon />}
            variant="outlined"
            color="secondary"
            type="submit">
            {isMobile ? 'Complete' : 'Mark as complete'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}
BuyOrderUpdateDialog.propTypes = {
  market: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
  onRemove: PropTypes.func,
  onSuccess: PropTypes.func,
}
BuyOrderUpdateDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
  onRemove: () => {},
  onSuccess: () => {},
}
