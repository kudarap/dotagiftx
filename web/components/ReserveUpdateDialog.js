import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import CircularProgress from '@material-ui/core/CircularProgress'
import DeliveredIcon from '@material-ui/icons/AssignmentTurnedIn'
import CancelIcon from '@material-ui/icons/Cancel'
import { myMarket } from '@/service/api'
import Button from '@/components/Button'
import ItemImage from '@/components/ItemImage'
import { amount, dateCalendar } from '@/lib/format'
import DialogCloseButton from '@/components/DialogCloseButton'
import {
  MARKET_STATUS_CANCELLED,
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_SOLD,
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

export default function ReserveUpdateDialog(props) {
  const classes = useStyles()

  const { market, open, onClose } = props

  const [error, setError] = React.useState('')
  const [loading, setLoading] = React.useState(false)

  const handleClose = () => {
    setError('')
    setLoading(false)
    onClose()
  }

  const router = useRouter()

  const marketUpdate = payload => {
    if (loading) {
      return
    }

    setLoading(true)
    setError(null)
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, payload)
        handleClose()
        router.push(`/history#${MARKET_STATUS_MAP_TEXT[payload.status].toLowerCase()}`)
      } catch (e) {
        setError(`Error: ${e.message}`)
      }

      setLoading(false)
    })()
  }

  const handleCancelClick = () => {
    marketUpdate({
      status: MARKET_STATUS_CANCELLED,
    })
  }

  const onFormSubmit = evt => {
    evt.preventDefault()

    marketUpdate({
      status: MARKET_STATUS_SOLD,
    })
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
          Update Reservation
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
                <Typography color="textSecondary" component="span">
                  {`Price: `}
                </Typography>
                {amount(market.price, market.currency)}
                <br />
                <Typography color="textSecondary" component="span">
                  {`Reserved Date: `}
                </Typography>
                {dateCalendar(market.updated_at)}
                {market.notes && (
                  <>
                    <br />
                    <Typography color="textSecondary" component="span">
                      {`Reserve Notes: `}
                    </Typography>
                    {market.notes}
                  </>
                )}
              </Typography>
            </Typography>
          </div>
        </DialogContent>
        {error && (
          <Typography color="error" align="center" variant="body2">
            {error}
          </Typography>
        )}
        <DialogActions>
          <Button disabled={loading} startIcon={<CancelIcon />} onClick={handleCancelClick}>
            Cancel Reservation
          </Button>
          <Button
            startIcon={
              loading ? <CircularProgress size={22} color="secondary" /> : <DeliveredIcon />
            }
            variant="outlined"
            color="secondary"
            type="submit">
            Item Delivered to Buyer
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
}
ReserveUpdateDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
}
