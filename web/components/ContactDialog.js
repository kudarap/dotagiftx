import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
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
import AppContext from '@/components/AppContext'
import MarketNotes from '@/components/MarketNotes'
import ProfileCard from '@/components/ProfileCard'

const marketSummaryFilter = {}

export default function ContactDialog(props) {
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
          <ProfileCard user={market.user} marketSummary={marketSummary} loading={loading}>
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
          </ProfileCard>

          <Typography variant="body2" color="textSecondary" component="div">
            <br />
            Guides for buying Giftables
            <ul>
              <li>
                Always check the item or set availability on seller&apos;s Dota 2 {` `}
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
                Dota 2 Giftables transaction only viable if the two steam user parties have been
                friends for 30 days.
              </li>
              <li>
                As Giftables involves a party having to go first, please always check seller&apos;s
                reputation through&nbsp;
                <Link
                  style={{ textDecoration: 'underline' }}
                  href={`${STEAMREP_PROFILE_BASE_URL}/${market.user.steam_id}`}
                  target="_blank"
                  rel="noreferrer noopener">
                  SteamRep
                </Link>
                &nbsp;and{' '}
                <Link
                  style={{ textDecoration: 'underline' }}
                  href={`/profiles/${market.user.steam_id}/delivered`}>
                  transaction history
                </Link>
                .
              </li>

              <li>
                If you need a middleman, I only suggest you get{' '}
                <Link href="/middlemen" target="_blank" color="secondary">
                  Middleman here
                </Link>
                .
              </li>

              {/* <li> */}
              {/*  Official SteamRep middleman may assist in middle manning for the trade, or{' '} */}
              {/*  <Link */}
              {/*    style={{ textDecoration: 'underline' }} */}
              {/*    href="https://www.reddit.com/r/dota2trade/" */}
              {/*    target="_blank" */}
              {/*    rel="noreferrer noopener"> */}
              {/*    r/Dota2Trade */}
              {/*  </Link>{' '} */}
              {/*  mod may assist as well in this. */}
              {/* </li> */}
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
