import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import Button from '@/components/Button'
import ItemImage from '@/components/ItemImage'
import { amount, dateCalendar } from '@/lib/format'
import DialogCloseButton from '@/components/DialogCloseButton'

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

  if (!market) {
    return null
  }

  return (
    <div>
      <Dialog
        fullWidth
        open={open}
        onClose={onClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
        <DialogTitle id="alert-dialog-title">
          Update Listing
          <DialogCloseButton onClick={onClose} />
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
                <strong>Live</strong>
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
                <br />
                <Typography color="textSecondary" component="span">
                  {`Notes: `}
                </Typography>
                {market.notes}
              </Typography>
            </Typography>
          </div>
        </DialogContent>
        <DialogActions>
          <Button>Remove</Button>
          <Button variant="outlined" color="secondary">
            Reserve to Buyer
          </Button>
        </DialogActions>
      </Dialog>
    </div>
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
