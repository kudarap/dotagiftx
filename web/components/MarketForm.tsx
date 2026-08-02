import React, { useContext, useEffect } from 'react'
import { makeStyles } from 'tss-react/mui'
import startsWith from 'lodash/startsWith'
import Paper from '@mui/material/Paper'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'
import CircularProgress from '@mui/material/CircularProgress'
import SubmitIcon from '@mui/icons-material/Check'
import Alert from '@mui/material/Alert'
import Box from '@mui/material/Box'
import { catalog, myMarket, myProfile } from '@/service/api'
import { APP_NAME } from '@/constants/strings'
import { USER_SUBSCRIPTION_MAP_COLOR } from '@/constants/user'
import { itemRarityColorMap } from '@/constants/palette'
import * as format from '@/lib/format'
import * as url from '@/lib/url'
import Button from '@/components/Button'
import ItemAutoComplete from '@/components/ItemAutoComplete'
import ItemImage from '@/components/ItemImage'
import Link from '@/components/Link'
import { MARKET_NOTES_MAX_LEN, MARKET_QTY_LIMIT } from '@/constants/market'
import { VERIFIED_INVENTORY_VERIFIED, VERIFIED_DELIVERY_MAP_ICON } from '@/constants/verified'
import AppContext from '@/components/AppContext'
import ReSellInput from './ReSellerInput'
import RefresherOrbBoon from './RefresherOrbBoon'
import type { TextFieldProps } from '@mui/material/TextField'
import type { Item } from '@/lib/types'

const defaultItem: Partial<Item> = {
  id: '',
}
const defaultPayload = {
  item_id: '',
  price: '',
  qty: 1,
  notes: '',
}

const steamCommunityBaseURL = 'https://steamcommunity.com'

const checkMarketPayload = (payload: {
  item_id?: string
  price: number
  quantity?: number
  notes: string
  seller_steam_id?: string
}) => {
  if (!payload.item_id) {
    return 'Item reference should be valid'
  }

  if (Number(payload.price) <= 0) {
    return 'Price must be atleast 0.01 USD'
  }

  if (Number(payload.quantity) > MARKET_QTY_LIMIT) {
    return `Quantity limit ${MARKET_QTY_LIMIT} per post`
  }

  const notesLen = String(payload.notes).length
  if (notesLen > MARKET_NOTES_MAX_LEN) {
    return `Notes max length limit reached ${notesLen}/${MARKET_NOTES_MAX_LEN}`
  }

  if (payload.seller_steam_id) {
    if (!url.isValid(payload.seller_steam_id)) {
      return 'Steam Profile is not a valid URL.'
    }
    if (!startsWith(payload.seller_steam_id, steamCommunityBaseURL, 0)) {
      return `Steam Profile should start with ${steamCommunityBaseURL}`
    }
  }

  return null
}

export default function MarketForm() {
  const { isLoggedIn } = useContext(AppContext)

  const [item, setItem] = React.useState<Partial<Item>>(defaultItem)
  const [payload, setPayload] = React.useState(defaultPayload)
  const [newMarketID, setNewMarketID] = React.useState<string | null>(null)
  const [partnerSteamID, setPartnerSteamID] = React.useState<string | null>(null)
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
      setBoons(user?.boons || [])
    })()
  }, [isLoggedIn])

  const handleItemSelect = (val: Partial<Item>) => {
    // Reset values when item is selected
    const newPayload = { ...defaultPayload }
    setPayload(newPayload)
    setNewMarketID(null)
    setError(null)

    setItem(val)
    // get item starting price
    if (val.slug) {
      setPayload({ ...newPayload, item_id: val.slug })
      catalog(val.slug)
        .then(res => {
          setItem(res as Item)
        })
        .catch(e => {
          console.log('error getting catalog info', (e as Error).message)
        })
    }
  }

  const handleSubmit = (evt: React.FormEvent) => {
    evt.preventDefault()

    // format and validate payload
    const quantity = Number(payload.qty)
    const newMarket: {
      item_id: string
      price: number
      notes: string
      seller_steam_id?: string
    } = {
      item_id: payload.item_id,
      price: Number(payload.price),
      notes: String(payload.notes).trim(),
    }

    // experimental fields
    if (partnerSteamID) {
      newMarket.seller_steam_id = String(partnerSteamID).trim()
    }

    const err = checkMarketPayload({ ...newMarket, quantity })
    if (err) {
      setError(`Error: ${err}`)
      return
    }

    setLoading(true)
    setError(null)
    setNewMarketID(null)
    ;(async () => {
      try {
        let res: { id: string } | undefined
        for (let i = 0; i < quantity; i++) {
          res = (await myMarket.POST(newMarket)) as { id: string }
        }

        // redirect to user listings
        if (res) {
          setNewMarketID(res.id)
          // setError('Item posted successfully! You will be redirected to your item listings.')
          // setTimeout(() => {
          //   router.push('/my-listings')
          // }, 3000)
        }
      } catch (e) {
        // special case reword error
        let m
        switch ((e as Error).message) {
          case 'market ask should be higher than highest bid price':
            m = 'sell price should be higher than current buy order price'
            break
          case 'user has been reported for scam incident':
            m = 'your account has been disabled due to scam report. please contact site admin'
            break
          default:
            m = (e as Error).message
        }

        setError(`Error: ${m}`)
      }

      setLoading(false)
    })()
  }

  const itemSelectEl = React.useRef<HTMLInputElement>(null)
  const handleFormReset = () => {
    setItem(defaultItem)
    setPayload(defaultPayload)
    setNewMarketID(null)
    setPartnerSteamID(null)
    setError(null)

    if (itemSelectEl.current) {
      const inputEl = itemSelectEl.current
      // inputEl.focus()
      inputEl.select()
    }
  }

  const handlePriceChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setPayload({ ...payload, price: e.target.value })

  const handleQtyChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const qty = e.target.value
    setPayload({ ...payload, qty: Number(qty) })
  }

  const subscribersColor = USER_SUBSCRIPTION_MAP_COLOR[subscription || 0]

  return (
    <>
      {!isLoggedIn && (
        <>
          <Alert severity="warning">
            You must be signed in to post an item — <Link href="/login">Sign in now</Link>
          </Alert>
          <br />
        </>
      )}

      <Paper
        component="form"
        sx={theme => ({
          maxWidth: theme.breakpoints.values.sm,
          margin: '0 auto',
          padding: theme.spacing(2),

          transition: `box-shadow .5s ease-in-out, border .2s`,
          borderTop: subscribersColor ? `5px solid ${subscribersColor}` : undefined,
          boxShadow: subscribersColor ? `0 0 15px ${subscribersColor}` : undefined,
        })}
        onSubmit={handleSubmit}>
        <Typography variant="h5" component="h1">
          Post your item on {APP_NAME}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Only verified ({VERIFIED_DELIVERY_MAP_ICON[VERIFIED_INVENTORY_VERIFIED]}) items from your
          inventory will be listed on Item page. All your posts will still be visible on your
          profile.
        </Typography>
        <br />

        <ItemAutoComplete
          required
          ref={itemSelectEl}
          onSelect={handleItemSelect}
          disabled={loading || !isLoggedIn}
        />
        <br />

        {/* Selected item preview */}
        {item.id && (
          <Box sx={{ display: 'flex', mb: 2 }}>
            <ItemImage
              style={{
                width: 150,
                height: 100,
              }}
              image={item.image || ''}
              width={150}
              height={100}
              rarity={item.rarity || ''}
              title={item.name || ''}
            />
            <Box sx={{ ml: 1 }}>
              <Typography variant="body2" color="textSecondary">
                Origin:{' '}
                <Typography variant="body2" color="textPrimary" component="span">
                  {item.origin}
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
                    color: itemRarityColorMap[item.rarity || ''],
                  }}>
                  {item.rarity}
                </Typography>
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Hero:{' '}
                <Typography variant="body2" color="textPrimary" component="span">
                  {item.hero}
                </Typography>
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Starting at:{' '}
                <Link href={`/${item.slug}`}>
                  {item.lowest_ask ? format.amount(item.lowest_ask, 'USD') : 'no offers yet'}
                </Link>
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Request to buy at:{' '}
                <Link
                  sx={{
                    color: 'accent.main',
                  }}
                  href={`/${item.slug}/buyorders`}>
                  {item.highest_bid ? format.amount(item.highest_bid, 'USD') : 'no orders yet'}
                </Link>
              </Typography>
            </Box>
          </Box>
        )}

        {boons && boons.indexOf('SHOPKEEPERS_CONTRACT') !== -1 && (
          <ReSellInput
            variant="outlined"
            fullWidth
            color="secondary"
            label="Seller Profile URL"
            placeholder="https://steamcommunity.com/..."
            value={partnerSteamID || ''}
            onInput={((e: React.FormEvent<HTMLInputElement>) =>
              setPartnerSteamID(e.currentTarget.value)) as unknown as NonNullable<TextFieldProps['onInput']>}
            disabled={loading || !isLoggedIn || Boolean(newMarketID)}
          />
        )}

        <div>
          <TextField
            variant="outlined"
            required
            color="secondary"
            label="Price"
            placeholder="1.00"
            type="number"
            helperText="Price value will be on USD."
            style={{ width: '69%' }}
            value={payload.price}
            onInput={handlePriceChange as unknown as NonNullable<TextFieldProps['onInput']>}
            onChange={handlePriceChange}
            onBlur={(e: React.FocusEvent<HTMLInputElement>) => {
              const price = format.amount(e.target.value)
              setPayload({ ...payload, price })
            }}
            disabled={loading || !isLoggedIn || Boolean(newMarketID)}
          />
          <TextField
            variant="outlined"
            color="secondary"
            label="Qty"
            type="number"
            value={payload.qty}
            style={{ width: '30%', marginLeft: '1%' }}
            onInput={handleQtyChange as unknown as NonNullable<TextFieldProps['onInput']>}
            onChange={handleQtyChange}
            onBlur={(e: React.FocusEvent<HTMLInputElement>) => {
              let qty = Number(e.target.value)
              if (qty < 1) {
                qty = 1
              }
              setPayload({ ...payload, qty })
            }}
            disabled={loading || !isLoggedIn || Boolean(newMarketID)}
          />
        </div>
        <br />

        <TextField
          variant="outlined"
          fullWidth
          color="secondary"
          label="Notes"
          value={payload.notes}
          helperText="Keep it short, This will be displayed when they check your offer."
          onInput={((e: React.FormEvent<HTMLInputElement>) =>
            setPayload({ ...payload, notes: e.currentTarget.value })) as unknown as NonNullable<TextFieldProps['onInput']>}
          disabled={loading || !isLoggedIn || Boolean(newMarketID)}
        />
        <br />
        <br />

        {!newMarketID && (
          <Button
            variant="contained"
            fullWidth
            type="submit"
            size="large"
            disabled={loading || !isLoggedIn || Boolean(newMarketID)}
            startIcon={loading ? <CircularProgress size={22} /> : <SubmitIcon />}>
            Post Item
          </Button>
        )}

        <div style={{ marginTop: 1 }}>
          {newMarketID && (
            <Alert
              severity="success"
              variant="filled"
              sx={{ color: 'primary.main' }}
              action={
                <Button color="inherit" size="small" onClick={handleFormReset}>
                  Post More
                </Button>
              }>
              Item posted successfully! Check your{' '}
              <Link style={{ textDecoration: 'underline' }} href="/my-listings">
                Item Listings
              </Link>
              .
            </Alert>
          )}
        </div>

        <RefresherOrbBoon boons={boons} />

        {error && (
          <Typography sx={{ mt: 1 }} align="center" variant="body2" color="error">
            {error}
          </Typography>
        )}

        <br />

        <Typography variant="body2" color="textSecondary" component="div">
          <strong>Guides for selling Giftables</strong>
          <ul>
            <li>Please make sure your item exist in your inventory.</li>
            <li>
              Dota 2 Giftables transaction only viable if the two steam user parties have been
              friends for 30 days.
            </li>
            <li>
              Please be clear in your terms and price. If the price is variable and subject to
              change, make a new post and remove the old one.
            </li>
            <li>
              Payment agreements will be done between you and the buyer. This website does not
              accept or integrate any payment service.
            </li>
          </ul>
        </Typography>
      </Paper>
    </>
  )
}
