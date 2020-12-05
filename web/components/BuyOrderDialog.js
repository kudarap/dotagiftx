import React, { useContext, useState } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import TextField from '@material-ui/core/TextField'
import Typography from '@material-ui/core/Typography'
import DialogTitle from '@material-ui/core/DialogTitle'
import DialogContent from '@material-ui/core/DialogContent'
import CircularProgress from '@material-ui/core/CircularProgress'
import SubmitIcon from '@material-ui/icons/Check'
import * as format from '@/lib/format'
import Link from '@/components/Link'
import ItemImage from '@/components/ItemImage'
import BidButton from '@/components/BidButton'
import AppContext from '@/components/AppContext'
import { MARKET_TYPE_BID } from '@/constants/market'
import { itemRarityColorMap } from '@/constants/palette'
import DialogCloseButton from '@/components/DialogCloseButton'

const useStyles = makeStyles(theme => ({
  details: {
    [theme.breakpoints.down('xs')]: {
      textAlign: 'center',
      display: 'block',
    },
    display: 'inline-flex',
  },
  profileName: {
    [theme.breakpoints.down('xs')]: {
      fontSize: theme.typography.h6.fontSize,
    },
  },
  avatar: {
    [theme.breakpoints.down('xs')]: {
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

const checkPayload = payload => {
  if (Number(payload.price) <= 0) {
    return 'Price must be atleast 0.01 USD'
  }

  return null
}

export default function BuyOrderDialog(props) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const { catalog, open, onClose } = props

  const [price, setPrice] = useState('')
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(false)

  const handleSubmit = e => {
    e.preventDefault()

    // format and validate payload
    const payload = {
      type: MARKET_TYPE_BID,
      item_id: catalog.item_id,
      price: Number(price),
    }
    const err = checkPayload(payload)
    if (err) {
      setError(`Error: ${err}`)
      return
    }

    setLoading(true)
    alert('market bid ready to fly')
  }

  const handlePriceChange = e => setPrice(e.target.value)

  return (
    <Dialog
      component="form"
      fullScreen={isMobile}
      open={open}
      onClose={onClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      <form onSubmit={handleSubmit}>
        <DialogTitle id="alert-dialog-title">
          Buy - {catalog.name}
          <DialogCloseButton onClick={onClose} />
        </DialogTitle>
        <DialogContent>
          <div>
            <ItemImage
              className={classes.itemImage}
              image={catalog.image}
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
                  color: itemRarityColorMap[catalog.rarity],
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
              Starting at:{' '}
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
                {format.amount(catalog.highest_bid, 'USD')}
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
              disabled={loading}
              value={price}
              onInput={handlePriceChange}
              onChange={handlePriceChange}
              onBlur={e => {
                const p = format.amount(e.target.value)
                setPrice(p)
              }}
            />
          </div>
          <br />

          <BidButton
            fullWidth
            disableUnderline
            size="large"
            type="submit"
            variant="outlined"
            target="_blank"
            rel="noreferrer noopener"
            disabled={loading}
            startIcon={loading ? <CircularProgress size={22} color="inherit" /> : <SubmitIcon />}>
            Place buy order
          </BidButton>

          {error && (
            <Typography align="center" variant="body2" color="error">
              {error}
            </Typography>
          )}

          <Typography variant="body2" color="textSecondary">
            <br />
            Placing buy order on Giftables
            <ul>
              <li>
                As giftables involves a party having to go first, please always check seller&apos;s
                reputation through&nbsp;
                <Link
                  style={{ textDecoration: 'underline' }}
                  href="https://steamrep.com"
                  target="_blank"
                  rel="noreferrer noopener">
                  SteamRep
                </Link>
                .
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
      </form>
    </Dialog>
  )
}
BuyOrderDialog.propTypes = {
  catalog: PropTypes.object.isRequired,
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
BuyOrderDialog.defaultProps = {
  open: false,
  onClose: () => {},
}
