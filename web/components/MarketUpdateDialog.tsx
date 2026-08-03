import React, { useContext } from 'react'
import moment from 'moment'
import startsWith from 'lodash/startsWith'
import { makeStyles } from 'tss-react/mui'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Typography from '@mui/material/Typography'
import TextField from '@mui/material/TextField'
import CircularProgress from '@mui/material/CircularProgress'
import ReserveIcon from '@mui/icons-material/EventAvailable'
import RemoveIcon from '@mui/icons-material/Delete'
import * as url from '@/lib/url'
import { amount, dateTime } from '@/lib/format'
import { myMarket } from '@/service/api'
import Button from '@/components/Button'
import Link from '@/components/Link'
import DialogCloseButton from '@/components/DialogCloseButton'
import {
  MARKET_STATUS_MAP_COLOR,
  MARKET_STATUS_MAP_TEXT,
  MARKET_STATUS_REMOVED,
  MARKET_STATUS_RESERVED,
} from '@/constants/market'
import AppContext from '@/components/AppContext'
import ItemImageDialog from '@/components/ItemImageDialog'
import type { TextFieldProps } from '@mui/material/TextField'
import type { Market } from '@/lib/types'

const useStyles = makeStyles()(theme => ({
  details: {
    [theme.breakpoints.down('sm')]: {
      display: 'block',
    },
    display: 'inline-flex',
  },
}))

const steamCommunityBaseURL = 'https://steamcommunity.com'
const steamProfileBaseURL = `${steamCommunityBaseURL}/profiles/`

interface UpdateDialogProps {
  market: Market | null
  open: boolean
  onClose: () => void
  onRemove?: () => void
  onSuccess?: () => void
}

export default function MarketUpdateDialog({
  market,
  open,
  onClose,
  onRemove,
  onSuccess,
}: UpdateDialogProps) {
  const { classes } = useStyles()
  const { isMobile } = useContext(AppContext)

  const [steamProfileURL, setSteamProfileURL] = React.useState('')
  const [notes, setNotes] = React.useState('')
  const [error, setError] = React.useState('')
  const [loading, setLoading] = React.useState(false)
  const [loadingRemove, setLoadingRemove] = React.useState(false)

  const handleClose = () => {
    setSteamProfileURL('')
    setNotes('')
    setError('')
    setLoading(false)
    setLoadingRemove(false)
    onClose()
  }

  const handleRemove = () => {
    if (onRemove) {
      onRemove()
    }
    handleClose()
  }

  const handleRemoveClick = () => {
    if (!market) {
      return
    }
    setLoadingRemove(true)
    setError('')
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, { status: MARKET_STATUS_REMOVED })
        handleRemove()
      } catch (e) {
        setError(`Error: ${(e as Error).message}`)
      }

      setLoadingRemove(false)
    })()
  }

  const onFormSubmit = (evt: React.FormEvent) => {
    evt.preventDefault()

    if (!market) {
      return
    }

    // if (loading || notes.trim() === '') {
    //   return
    // }
    if (!url.isValid(steamProfileURL)) {
      setError('Steam Profile is not a valid URL.')
      return
    }
    if (!startsWith(steamProfileURL, steamCommunityBaseURL, 0)) {
      setError(`Steam Profile should start with ${steamCommunityBaseURL}`)
      return
    }

    const payload = {
      status: MARKET_STATUS_RESERVED,
      partner_steam_id: steamProfileURL,
      notes,
    }

    setLoading(true)
    setError('')
    ;(async () => {
      try {
        await myMarket.PATCH(market.id, payload)
        handleClose()
        if (onSuccess) {
          onSuccess()
        }
      } catch (e) {
        setError(`Error: ${(e as Error).message}`)
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
      aria-describedby="alert-dialog-description"
    >
      <form onSubmit={onFormSubmit}>
        <DialogTitle id="alert-dialog-title">
          Update Listing
          <DialogCloseButton onClick={handleClose} />
        </DialogTitle>
        <DialogContent>
          <div className={classes.details}>
            <ItemImageDialog item={market.item as NonNullable<Market['item']>} />

            <Typography component="h1">
              <Typography variant="h6" component={Link} href={`/${market.item?.slug}`}>
                {market.item?.name}
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
                      rel="noreferrer noopener"
                    >
                      {market.seller_steam_id}
                    </Link>
                    <br />
                  </div>
                )}
                <Typography color="textSecondary" component="span">
                  {`Price: `}
                </Typography>
                {amount(market.price || 0, market.currency)}
                <br />
                <Typography color="textSecondary" component="span">
                  {`Listed: `}
                </Typography>
                {dateTime(market.updated_at || '')}
                {market.notes && (
                  <>
                    <br />
                    <Typography color="textSecondary" component="span">
                      {`Notes: `}
                    </Typography>
                    <Typography component="ul" variant="body2" style={{ marginTop: 0 }}>
                      {market.notes.split('\n').map(s => (
                        <li key={s}>{s}</li>
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
              placeholder="https://steamcommunity.com/..."
              value={steamProfileURL}
              onInput={
                ((e: React.FormEvent<HTMLInputElement>) =>
                  setSteamProfileURL(e.currentTarget.value)) as unknown as NonNullable<
                  TextFieldProps['onInput']
                >
              }
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
              placeholder={`${moment().add(30, 'days').format('MMM D')} - $1 deposit`}
              value={notes}
              onInput={
                ((e: React.FormEvent<HTMLInputElement>) =>
                  setNotes(e.currentTarget.value)) as unknown as NonNullable<
                  TextFieldProps['onInput']
                >
              }
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
            disabled={loadingRemove}
            startIcon={loadingRemove ? <CircularProgress size={22} /> : <RemoveIcon />}
            onClick={handleRemoveClick}
            variant="outlined"
          >
            Remove Listing
          </Button>
          <Button
            disabled={loading}
            startIcon={loading ? <CircularProgress size={22} color="secondary" /> : <ReserveIcon />}
            variant="outlined"
            color="secondary"
            type="submit"
          >
            Reserve to Buyer
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}
