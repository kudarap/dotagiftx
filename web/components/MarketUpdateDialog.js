import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
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

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  media: {
    [theme.breakpoints.down('xs')]: {
      margin: '0 auto !important',
    },
    width: 150,
    height: 100,
    marginRight: theme.spacing(1.5),
    marginBottom: theme.spacing(1.5),
  },
}))

export default function MarketUpdateDialog(props) {
  const classes = useStyles()

  const { market, open, onClose } = props

  const [notes, setNotes] = React.useState('')
  const [error, setError] = React.useState('')
  const [loading, setLoading] = React.useState(false)

  const handleClose = () => {
    setNotes('')
    setError('')
    setLoading(false)
    onClose()
  }

  const router = useRouter()
  const handleRemoveClick = () => {
    setLoading(true)
    setError(null)
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, { status: MARKET_STATUS_REMOVED })
        handleClose()
        router.reload()
      } catch (e) {
        setError(`Error: ${e.message}`)
      }

      setLoading(false)
    })()
  }

  const onFormSubmit = evt => {
    evt.preventDefault()

    if (loading || notes.trim() === '') {
      return
    }

    const payload = {
      status: MARKET_STATUS_RESERVED,
      notes,
    }

    setLoading(true)
    setError(null)
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, payload)
        handleClose()
        router.push('/reservations')
      } catch (e) {
        setError(`Error: ${e.message}`)
      }

      setLoading(false)
    })()
  }

  if (!market) {
    return null
  }

  return (
    <Dialog
      fullWidth
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
            <ItemImage
              className={classes.media}
              image={`/300x170/${market.item.image}`}
              title={market.item.name}
              rarity={market.item.rarity}
            />

            <Typography component="h1">
              <Typography
                component="p"
                variant="h6"
                component={Link}
                href="/item/[slug]"
                as={`/item/${market.item.slug}`}>
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
                    {market.notes}
                  </>
                )}
              </Typography>
            </Typography>
          </div>
          <div>
            <TextField
              disabled={loading}
              fullWidth
              required
              color="secondary"
              variant="outlined"
              label="Reserved Notes"
              helperText="Helpful info like Buyer's avatar name, Steam profile url, or Delivery date"
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
          <Button disabled={loading} startIcon={<RemoveIcon />} onClick={handleRemoveClick}>
            Remove listing
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
}
MarketUpdateDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
}
