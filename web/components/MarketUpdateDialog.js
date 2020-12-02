import React, { useContext } from 'react'
import PropTypes from 'prop-types'
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
import { myMarket } from '@/service/api'
import { amount, dateCalendar } from '@/lib/format'
import Button from '@/components/Button'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import DialogCloseButton from '@/components/DialogCloseButton'
import {
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_REMOVED,
  MARKET_STATUS_RESERVED,
} from '@/constants/market'
import AppContext from '@/components/AppContext'

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      display: 'block',
    },
    display: 'inline-flex',
  },
  media: {
    [theme.breakpoints.down('xs')]: {
      margin: '0 auto !important',
      width: 300,
      height: 170,
    },
    width: 165,
    height: 110,
    margin: theme.spacing(0, 1.5, 1.5, 0),
  },
}))

export default function MarketUpdateDialog(props) {
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

    if (loading || notes.trim() === '') {
      return
    }

    const payload = {
      status: MARKET_STATUS_RESERVED,
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
          Update Listing
          <DialogCloseButton onClick={handleClose} />
        </DialogTitle>
        <DialogContent>
          <div className={classes.details}>
            {isMobile ? (
              <ItemImage
                className={classes.media}
                image={market.item.image}
                width={300}
                height={170}
                title={market.item.name}
                rarity={market.item.rarity}
              />
            ) : (
              <ItemImage
                className={classes.media}
                image={market.item.image}
                width={165}
                height={110}
                title={market.item.name}
                rarity={market.item.rarity}
              />
            )}

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
                {dateCalendar(market.updated_at)}
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
              label="Buyer's Steam profile URL"
              placeholder="https://steamcommunity.com/profiles/..."
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
              label="Reservation Notes"
              helperText="Delivery date and deposit details"
              placeholder="Jan 2 with $1 deposit"
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
            Remove Listing
          </Button>
          <Button
            startIcon={loading ? <CircularProgress size={22} color="secondary" /> : <ReserveIcon />}
            variant="outlined"
            color="secondary"
            type="submit">
            Reserve to Buyer
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}
MarketUpdateDialog.propTypes = {
  market: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
  onRemove: PropTypes.func,
  onSuccess: PropTypes.func,
}
MarketUpdateDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
  onRemove: () => {},
  onSuccess: () => {},
}
