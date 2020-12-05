import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import { Avatar, TextField } from '@material-ui/core'
import * as format from '@/lib/format'
import ChipLink from '@/components/ChipLink'
import { STEAM_PROFILE_BASE_URL, STEAMREP_PROFILE_BASE_URL } from '@/constants/strings'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import ItemImage, { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import BidButton from '@/components/BidButton'
import { itemRarityColorMap } from '@/constants/palette'

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
}))

export default function BuyOrderDialog(props) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const { catalog, open, onClose } = props

  return (
    <div>
      <Dialog
        fullWidth
        fullScreen={isMobile}
        open={open}
        onClose={onClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
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

          <TextField
            variant="outlined"
            required
            color="secondary"
            label="Price"
            placeholder="1.00"
            type="number"
            helperText="Price you want to pay in USD."
          />

          <Typography variant="body2" color="textSecondary">
            <br />
            Guides for placing buy order on Giftables
            <ul>
              {/*<li>*/}
              {/*  Dota 2 giftables transaction only viable if the two steam user parties have been*/}
              {/*  friends for 30 days.*/}
              {/*</li>*/}
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
              {/*<li>*/}
              {/*  Official SteamRep middleman may assist in middle manning for the trade, or{' '}*/}
              {/*  <Link*/}
              {/*    style={{ textDecoration: 'underline' }}*/}
              {/*    href="https://www.reddit.com/r/dota2trade/"*/}
              {/*    target="_blank"*/}
              {/*    rel="noreferrer noopener">*/}
              {/*    r/Dota2Trade*/}
              {/*  </Link>{' '}*/}
              {/*  mod may assist as well in this.*/}
              {/*</li>*/}
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
        <DialogActions>
          {/* <Button component="a">Buyer Profile</Button> */}
          <BidButton variant="outlined" target="_blank" rel="noreferrer noopener" disableUnderline>
            Place buy order
          </BidButton>
        </DialogActions>
      </Dialog>
    </div>
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
