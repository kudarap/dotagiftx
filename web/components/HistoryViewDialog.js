import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import Dialog from '@material-ui/core/Dialog'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import { myMarket } from '@/service/api'
import ItemImage from '@/components/ItemImage'
import { amount, dateCalendar } from '@/lib/format'
import DialogCloseButton from '@/components/DialogCloseButton'
import {
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

export default function HistoryViewDialog(props) {
  const classes = useStyles()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))

  const { market, open, onClose } = props

  const [error, setError] = React.useState('')
  const [loading, setLoading] = React.useState(false)

  const handleClose = () => {
    setError('')
    setLoading(false)
    onClose()
  }

  const router = useRouter()

  const onFormSubmit = evt => {
    evt.preventDefault()

    if (loading) {
      return
    }

    const payload = {
      status: MARKET_STATUS_SOLD,
    }

    setLoading(true)
    setError(null)
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, payload)
        handleClose()
        router.push('/history')
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
      fullScreen={isMobile}
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      <DialogTitle id="alert-dialog-title">
        Listing Details
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
                {`Listed: `}
              </Typography>
              {dateCalendar(market.created_at)}
              <br />

              <Typography color="textSecondary" component="span">
                {`Updated: `}
              </Typography>
              {dateCalendar(market.updated_at)}
              <br />

              {market.notes && (
                <>
                  <Typography color="textSecondary" component="span">
                    {`Notes: `}
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
    </Dialog>
  )
}
HistoryViewDialog.propTypes = {
  market: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
HistoryViewDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
}
