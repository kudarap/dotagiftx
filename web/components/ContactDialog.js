import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import { Avatar } from '@material-ui/core'
import { statsMarketSummary } from '@/service/api'
import ChipLink from '@/components/ChipLink'
import {
  DOTABUFF_PROFILE_BASE_URL,
  STEAM_PROFILE_BASE_URL,
  STEAMREP_PROFILE_BASE_URL,
} from '@/constants/strings'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import MarketNotes from '@/components/MarketNotes'

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
}))

const marketSummaryFilter = {}

export default function ContactDialog(props) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const { market, open, onClose } = props

  const [loading, setLoading] = React.useState(true)
  const [marketSummary, setMarketSummary] = React.useState(null)
  React.useEffect(() => {
    if (!market) {
      return
    }

    ;(async () => {
      marketSummaryFilter.user_id = market.user.id
      try {
        const res = await statsMarketSummary(marketSummaryFilter)
        setMarketSummary(res)
      } catch (e) {
        console.log('error getting stats market summary', e.message)
      }
      setLoading(false)
    })()

    // eslint-disable-next-line consistent-return
    return () => {
      setMarketSummary(null)
    }
  }, [market])

  if (!market) {
    return null
  }

  const storeProfile = `/profiles/${market.user.steam_id}`
  const steamProfileURL = `${STEAM_PROFILE_BASE_URL}/${market.user.steam_id}`
  const dota2Inventory = `${steamProfileURL}/inventory#570`

  const reported = market.user.status && market.user.status === 300

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
          Contact Seller
          <DialogCloseButton onClick={onClose} />
        </DialogTitle>
        <DialogContent>
          <div
            className={classes.details}
            style={reported ? { backgroundColor: '#2d0000', padding: 10, width: '100%' } : null}>
            <a href={storeProfile} target="_blank" rel="noreferrer noopener">
              <Avatar className={classes.avatar} {...retinaSrcSet(market.user.avatar, 100, 100)} />
            </a>
            <Typography component="h1">
              <Typography
                className={classes.profileName}
                component="p"
                variant="h4"
                color={reported ? 'error' : ''}>
                {market.user.name}
              </Typography>
              {reported ? (
                <Typography color="error">
                  This user was reported over scam report and under investigation.
                </Typography>
              ) : null}
              <Typography variant="body2" component="span">
                <Link href={`/profiles/${market.user.steam_id}`}>
                  {!loading && marketSummary ? marketSummary.live : '--'} Items
                </Link>{' '}
                &middot;{' '}
                <Link href={`/profiles/${market.user.steam_id}/reserved`}>
                  {!loading && marketSummary ? marketSummary.reserved : '--'} Reserved
                </Link>{' '}
                &middot;{' '}
                <Link href={`/profiles/${market.user.steam_id}/delivered`}>
                  {!loading && marketSummary ? marketSummary.sold : '--'} Delivered
                </Link>
              </Typography>
              <br />
              <Typography gutterBottom>
                {/* <Typography color="textSecondary" component="span"> */}
                {/*  {`Links: `} */}
                {/* </Typography> */}
                {/* <ChipLink label="Steam Profile" href={steamProfileURL} /> */}
                {/* &nbsp; */}
                <ChipLink
                  label="SteamRep"
                  href={`${STEAMREP_PROFILE_BASE_URL}/${market.user.steam_id}`}
                />
                &nbsp;
                <ChipLink
                  label="Dotabuff"
                  href={`${DOTABUFF_PROFILE_BASE_URL}/${market.user.steam_id}`}
                />
                &nbsp;
                <ChipLink label="Steam Inventory" href={dota2Inventory} />
                {market.notes && <MarketNotes text={market.notes} />}
              </Typography>
            </Typography>
          </div>

          <Typography variant="body2" color="textSecondary" component="div">
            <br />
            Guides for buying Giftables
            <ul>
              <li>
                Always check the item/set availability on seller&apos;s Dota 2 {` `}
                <Link
                  style={{ textDecoration: 'underline' }}
                  href={dota2Inventory}
                  target="_blank"
                  rel="noreferrer noopener">
                  inventory
                </Link>
                .
              </li>
              <li>
                Dota 2 giftables transaction only viable if the two steam user parties have been
                friends for 30 days.
              </li>
              <li>
                As giftables involves a party having to go first, please always check seller&apos;s
                reputation through&nbsp;
                <Link
                  style={{ textDecoration: 'underline' }}
                  href={`${STEAMREP_PROFILE_BASE_URL}/${market.user.steam_id}`}
                  target="_blank"
                  rel="noreferrer noopener">
                  SteamRep
                </Link>
                .
              </li>

              <li>
                Official SteamRep middleman may assist in middle manning for the trade, or{' '}
                <Link
                  style={{ textDecoration: 'underline' }}
                  href="https://www.reddit.com/r/dota2trade/"
                  target="_blank"
                  rel="noreferrer noopener">
                  r/Dota2Trade
                </Link>{' '}
                mod may assist as well in this.
              </li>
            </ul>
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button component="a" href={storeProfile}>
            View Seller Items
          </Button>
          <Button
            color="secondary"
            variant="outlined"
            component={Link}
            disableUnderline
            target="_blank"
            rel="noreferrer noopener"
            href={steamProfileURL}>
            Check Steam Profile
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  )
}
ContactDialog.propTypes = {
  market: PropTypes.object,
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
ContactDialog.defaultProps = {
  market: null,
  open: false,
  onClose: () => {},
}
