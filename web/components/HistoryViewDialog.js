import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import ItemImage from '@/components/ItemImage'
import { amount, dateCalendar } from '@/lib/format'
import DialogCloseButton from '@/components/DialogCloseButton'
import { MARKET_STATUS_MAP_COLOR, MARKET_STATUS_MAP_TEXT } from '@/constants/market'
import AppContext from '@/components/AppContext'
import { TextField } from '@material-ui/core'

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

export default function HistoryViewDialog(props) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const { market, open, onClose } = props

  const handleClose = () => {
    onClose()
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
        Item Details
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
            style={{ marginTop: 16, marginBottom: 16 }}
            disabled
            fullWidth
            color="secondary"
            variant="outlined"
            label="Buyer's Steam profile URL"
            value={market.partner_steam_id}
          />
        </div>
      </DialogContent>
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
