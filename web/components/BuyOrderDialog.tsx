import React, { useContext, useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import Dialog from '@mui/material/Dialog'
import TextField from '@mui/material/TextField'
import type { TextFieldProps } from '@mui/material/TextField'
import Typography from '@mui/material/Typography'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import Paper from '@mui/material/Paper'
import CircularProgress from '@mui/material/CircularProgress'
import Alert from '@mui/material/Alert'
import SubmitIcon from '@mui/icons-material/Check'
import * as format from '@/lib/format'
import { myMarket, myProfile } from '@/service/api'
import Link from '@/components/Link'
import ItemImage from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import { MARKET_NOTES_MAX_LEN, MARKET_TYPE_BID } from '@/constants/market'
import { itemRarityColorMap } from '@/constants/palette'
import DialogCloseButton from '@/components/DialogCloseButton'
import Button from '@/components/Button'
import { USER_SUBSCRIPTION_MAP_COLOR } from '@/constants/user'
import RefresherOrbBoon from './RefresherOrbBoon'
import RefresherShardBoon from './RefresherShardBoon'
import type { Item } from '@/lib/types'

const useStyles = makeStyles()(theme => ({
  details: {
    [theme.breakpoints.down('sm')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  profileName: {
    [theme.breakpoints.down('sm')]: {
      fontSize: theme.typography.h6.fontSize,
    },
  },
  avatar: {
    [theme.breakpoints.down('sm')]: {
      margin: '0 auto',
    },
    width: 100,
    height: 100,
    marginRight: theme.spacing(1.5),
  },
  itemImage: {
    width: 150,
    height: 100,
    float: 'left',
    marginRight: theme.spacing(1),
  },
  spacer: {
    width: theme.spacing(1),
  },
}))

type Catalog = Item & {
  id: string
  origin?: string
  hero?: string
  bid_count?: number
  highest_bid?: number
  lowest_ask?: number
}

const checkPayload = (payload: { price: number; notes: string }) => {
  if (Number(payload.price) < 0.5) {
    return 'Price must be atleast 0.50 USD'
  }

  const notesLen = String(payload.notes).length
  if (notesLen > MARKET_NOTES_MAX_LEN) {
    return `Notes max length limit reached ${notesLen}/${MARKET_NOTES_MAX_LEN}`
  }

  return null
}

export default function BuyOrderDialog({
  catalog,
  open,
  onClose,
  onChange,
}: {
  catalog: Catalog
  open: boolean
  onClose: () => void
  onChange: () => void
}) {
  const { classes } = useStyles()
  const { isMobile, isLoggedIn } = useContext(AppContext)

  const [market, setMarket] = useState<unknown | null>(null)
  const [price, setPrice] = useState('')
  const [notes, setNotes] = useState('')
  const [error, setError] = React.useState<string | null>(null)
  const [loading, setLoading] = React.useState(false)

  const [subscription, setSubscription] = React.useState<number | null>(null)
  const [boons, setBoons] = React.useState<string[]>([])
  useEffect(() => {
    if (!isLoggedIn) {
      return
    }

    ;(async () => {
      const user = (await myProfile.GET(true)) as { subscription?: number; boons?: string[] }
      setSubscription(user?.subscription || null)
      setBoons([...(user?.boons || [])])
    })()
  }, [isLoggedIn])

  const handleSubmit = (evt: React.FormEvent) => {
    evt.preventDefault()

    // format and validate payload
    const buyOrder = {
      type: MARKET_TYPE_BID,
      item_id: catalog.id,
      price: Number(price),
      notes: String(notes).trim(),
    }
    const err = checkPayload(buyOrder)
    if (err) {
      setError(`Error: ${err}`)
      return
    }

    // submit buy order payload
    setError(null)
    setLoading(true)
    ;(async () => {
      try {
        const res = await myMarket.POST(buyOrder)
        setMarket(res)
        onChange()
      } catch (e) {
        // special case reword error
        let m = (e as Error).message
        if ((e as Error).message === 'market bid should be lower than lowest ask price') {
          m = 'buy order price should be lower than lowest offer price'
        }

        setError(`Error: ${m}`)
      }

      setLoading(false)
    })()
  }

  const router = useRouter()
  const handleClose = () => {
    onClose()
    setTimeout(() => {
      setError(null)
      setLoading(false)
      setNotes('')
      setPrice('')
      setMarket(null)
    }, 500)

    if (market) {
      // Forces to refresh buy order table
      router.push(`/${catalog.slug}/buyorders`)
    }
  }

  const handleDone = () => {
    handleClose()
  }

  const handlePriceChange = (evt: React.ChangeEvent<HTMLInputElement>) => setPrice(evt.target.value)
  const isInputDisabled = loading || Boolean(market) || !isLoggedIn

  const subscribersColor = subscription !== null ? USER_SUBSCRIPTION_MAP_COLOR[subscription] : undefined

  const refresherOrbOnly =
    (boons.indexOf('REFRESHER_ORB') !== -1 && boons.indexOf('REFRESHER_SHARD') !== -1) ||
    boons.indexOf('REFRESHER_ORB') !== -1

  return (
    <Dialog
      fullScreen={isMobile}
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      <Paper
        component="form"
        onSubmit={handleSubmit}
        sx={{
          transition: `box-shadow .5s ease-in-out, border .2s`,
          borderTop: subscribersColor ? `5px solid ${subscribersColor}` : 'none',
          boxShadow: subscribersColor ? `0 0 15px ${subscribersColor}` : 'none',
        }}>
        <DialogTitle id="alert-dialog-title">
          <DialogCloseButton onClick={handleClose} />
          Buy - {catalog.name}
        </DialogTitle>
        <DialogContent>
          {!isLoggedIn && (
            <>
              <Alert severity="warning">
                You must be signed in to place buy order — <Link href="/login">Sign in now</Link>
              </Alert>
              <br />
            </>
          )}

          <div>
            <ItemImage
              className={classes.itemImage}
              image={catalog.image || ''}
              width={150}
              height={100}
              rarity={catalog.rarity}
              title={catalog.name}
            />
            <Typography variant="body2" color="textSecondary">
              Origin:{' '}
              <Typography variant="body2" color="textPrimary" component="span">
                {catalog.origin}
              </Typography>
            </Typography>
            <Typography variant="body2" color="textSecondary">
              Rarity:{' '}
              <Typography
                variant="body2"
                color="textPrimary"
                component="span"
                style={{
                  textTransform: 'capitalize',
                  color: catalog.rarity ? itemRarityColorMap[catalog.rarity] : undefined,
                }}>
                {catalog.rarity}
              </Typography>
            </Typography>
            <Typography variant="body2" color="textSecondary">
              Hero:{' '}
              <Typography variant="body2" color="textPrimary" component="span">
                {catalog.hero}
              </Typography>
            </Typography>
            <Typography variant="body2" color="textSecondary">
              Lowest Price:{' '}
              <Link href={`/${catalog.slug}`}>
                {catalog.lowest_ask ? format.amount(catalog.lowest_ask, 'USD') : 'no offers yet'}
              </Link>
            </Typography>
            <br />
            <br />
          </div>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Typography color="textSecondary" style={{ marginTop: -22 }}>
              <Typography color="textPrimary" component="span">
                {catalog.bid_count}
              </Typography>{' '}
              requests to buy at{' '}
              <Typography color="textPrimary" component="span">
                {format.amount(catalog.highest_bid || 0, 'USD')}
              </Typography>{' '}
              or lower
            </Typography>
            <span className={classes.spacer} />
            <TextField
              required
              variant="outlined"
              color="secondary"
              label="Price"
              placeholder="1.00"
              type="number"
              helperText="Price you want to pay in USD."
              disabled={isInputDisabled}
              value={price}
              onInput={handlePriceChange as unknown as NonNullable<TextFieldProps['onInput']>}
              onChange={handlePriceChange}
              onBlur={(e: React.FocusEvent<HTMLInputElement>) => {
                const p = format.amount(e.target.value)
                setPrice(p)
              }}
            />
          </div>
          <br />
          <TextField
            disabled={isInputDisabled}
            fullWidth
            color="secondary"
            variant="outlined"
            label="Notes"
            helperText="Keep it short, This will be displayed when they check your order."
            value={notes}
            onInput={((e: React.FormEvent<HTMLInputElement>) =>
              setNotes(e.currentTarget.value)) as unknown as NonNullable<TextFieldProps['onInput']>}
          />
          <br />
          <br />
          {!market && (
            <Button
              fullWidth
              size="large"
              type="submit"
              variant="contained"
              disabled={isInputDisabled}
              startIcon={loading ? <CircularProgress size={22} color="inherit" /> : <SubmitIcon />}>
              Place Order
            </Button>
          )}

          {refresherOrbOnly ? (
            <RefresherOrbBoon boons={boons} />
          ) : (
            <RefresherShardBoon boons={boons} />
          )}

          <div style={{ marginTop: 2 }}>
            {Boolean(market) && (
              <Alert
                severity="success"
                variant="filled"
                sx={{ color: 'primary.main' }}
                action={
                  <Button color="inherit" size="small" onClick={handleDone}>
                    Done
                  </Button>
                }>
                Your buy order has been placed and now open for sellers. Check your{' '}
                <Link style={{ textDecoration: 'underline' }} href="/my-orders">
                  Buy orders
                </Link>
                .
              </Alert>
            )}
            {error && (
              <Typography align="center" variant="body2" color="error">
                {error}
              </Typography>
            )}
          </div>
          <Typography variant="body2" color="textSecondary" component="div">
            <br />
            Placing buy order on Giftables
            <ul>
              <li>
                Buyer profile details are{' '}
                <Typography component="span" variant="inherit" color="textPrimary">
                  only visible
                </Typography>{' '}
                to signed in website users.
              </li>
              <li>
                Payment agreements will be done between you and the seller. This website does not
                accept or integrate any payment service.
              </li>
              <li>
                Please kindly remove this buy order after use to prevent seller&apos;s contacting
                you.
              </li>
            </ul>
          </Typography>
        </DialogContent>
      </Paper>
    </Dialog>
  )
}
