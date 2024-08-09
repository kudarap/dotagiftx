import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Card from '@mui/material/Card'
import CardActions from '@mui/material/CardActions'
import CardContent from '@mui/material/CardContent'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import Popper from '@mui/material/Popper'
import { indigo } from '@mui/material/colors'
import { ThemeProvider } from '@mui/material/styles'
import { muiLightTheme } from '@/lib/theme'
import {
  VERIFIED_DELIVERY_MAP_LABEL,
  VERIFIED_DELIVERY_MAP_TEXT,
  VERIFIED_INVENTORY_MAP_LABEL,
  VERIFIED_INVENTORY_MAP_TEXT,
  VERIFIED_INVENTORY_PENDING,
} from '@/constants/verified'
import { dateFromNow } from '@/lib/format'
import Link from '@/components/Link'
import { MARKET_STATUS_SOLD } from '@/constants/market'

const useStyles = makeStyles()({
  root: {
    minWidth: 300,
    zIndex: 1,
  },
  poweredBy: {
    color: indigo[400],
  },
})

const assetModifier = asset => {
  let isGiftable = asset.gift_once ? 'Yes' : 'No'
  if (asset.type.startsWith('Immortal') && !asset.gift_once) {
    isGiftable = '?'
  }

  // Shows gift containing items.
  let displayName = asset.name
  let received = asset.date_received
  if (asset.name === 'Wrapped Gift') {
    displayName = asset.contains
    received = asset.name
  }

  return { ...asset, isGiftable, displayName, received }
}

const getInventoryURL = steamID => `https://steamcommunity.com/profiles/${steamID}/inventory/#570_2`

export default function VerifiedStatusCard({ market, ...other }) {
  const { classes } = useStyles()

  if (market === null) {
    return null
  }

  const { inventory, delivery } = market

  if (market.inventory_status == VERIFIED_INVENTORY_PENDING) {
    return <PendingCard data={market} />
  }

  let inventoryURL = getInventoryURL(market.user.steam_id)
  let source = inventory
  let mapLabel = VERIFIED_INVENTORY_MAP_LABEL
  let mapText = VERIFIED_INVENTORY_MAP_TEXT
  let isDelivery = false
  if (delivery) {
    isDelivery = true
    source = delivery
    mapText = VERIFIED_DELIVERY_MAP_TEXT
    mapLabel = VERIFIED_DELIVERY_MAP_LABEL
    inventoryURL = getInventoryURL(market.partner_steam_id)
  }
  if (!source) {
    if (market.resell) {
      return <ResellCard data={market} />
    }

    return null
  }

  const steamInvProfile =
    market.status === MARKET_STATUS_SOLD ? market.partner_steam_id : market.user.steam_id

  return (
    <CardX className={classes.root} {...other}>
      <CardContent>
        <Typography variant="h5" component="h2">
          {mapLabel[source.status]}
        </Typography>
        <Typography color="textSecondary" variant="body2" component="p">
          Last updated {dateFromNow(source.updated_at)}
        </Typography>
        <Typography component="p">{mapText[source.status]}</Typography>

        {source.steam_assets && (
          <>
            {!isDelivery && (
              <Typography variant="body2">
                <br />
                Found <strong>{source.bundle_count}</strong> bundle{source.bundle_count > 1 && 's'}
              </Typography>
            )}

            <Table className={classes.table} size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  {isDelivery ? (
                    <TableCell>From</TableCell>
                  ) : (
                    <TableCell align="center">Giftable</TableCell>
                  )}
                  {isDelivery ? <TableCell>Received</TableCell> : <TableCell>Qty</TableCell>}
                </TableRow>
              </TableHead>
              <TableBody>
                {source.steam_assets.map(assetModifier).map(asset => (
                  <TableRow key={asset.name}>
                    <TableCell component="th" scope="row">
                      <Link
                        color="secondary"
                        target="_blank"
                        rel="noreferrer noopener"
                        underline="none"
                        href={`${inventoryURL}_${asset.asset_id}`}>
                        <strong>{asset.displayName}</strong>
                      </Link>
                    </TableCell>
                    {isDelivery ? (
                      <TableCell>{asset.gift_from ? asset.gift_from : <ClearedGift />}</TableCell>
                    ) : (
                      <TableCell align="center">{asset.isGiftable}</TableCell>
                    )}
                    {isDelivery ? (
                      <TableCell>{asset.received ? asset.received : <ClearedGift />}</TableCell>
                    ) : (
                      <TableCell align="center">{asset.qty}</TableCell>
                    )}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </>
        )}
      </CardContent>
      <CardActions style={{ float: 'right' }}>
        <Link
          className={classes.poweredBy}
          variant="caption"
          target="_blank"
          rel="noreferrer noopener"
          underline="none"
          href={`https://steaminventory.org/?profile=${steamInvProfile}`}>
          Powered by <strong>SteamInventory.org</strong>
        </Link>
      </CardActions>
    </CardX>
  )
}

VerifiedStatusCard.propTypes = {
  market: PropTypes.object,
}
VerifiedStatusCard.defaultProps = {
  market: null,
}

function CardX(props) {
  return (
    <ThemeProvider theme={muiLightTheme}>
      <Card {...props} />
    </ThemeProvider>
  )
}

export function VerifiedStatusPopover({ market, ...other }) {
  return (
    <Popper
      style={{ marginTop: 2, zIndex: 1 }}
      placement="right-start"
      disablePortal={false}
      modifiers={[
        {
          name: 'flip',
          enabled: true,
          options: {
            altBoundary: true,
            rootBoundary: 'document',
            padding: 8,
          },
        },
        {
          name: 'preventOverflow',
          enabled: true,
          options: {
            altAxis: true,
            altBoundary: true,
            tether: true,
            rootBoundary: 'document',
            padding: 8,
          },
        },
      ]}
      {...other}>
      <VerifiedStatusCard market={market} onMouseLeave={other.onClose} />
    </Popper>
  )
}
VerifiedStatusPopover.propTypes = VerifiedStatusCard.propTypes
VerifiedStatusPopover.defaultProps = VerifiedStatusCard.defaultProps

function ClearedGift() {
  return (
    <Typography color="textSecondary" variant="body2" component="em">
      --
      {/* cleared */}
    </Typography>
  )
}

function ResellCard({ data }) {
  return (
    <CardX>
      <CardContent>
        <Typography variant="h5" component="h2">
          Item Resell
        </Typography>
        <Typography color="textSecondary" variant="body2" component="p">
          Verified by {data.user.name} {dateFromNow(data.created_at)}
        </Typography>
        <Typography component="p">
          Item manually verified by seller from partner&apos;s inventory
        </Typography>
      </CardContent>
    </CardX>
  )
}

function PendingCard({ data }) {
  return (
    <CardX>
      <CardContent>
        <Typography variant="h5" component="h2">
          Pending
        </Typography>
        <Typography color="textSecondary" variant="body2" component="p">
          Posted {dateFromNow(data.created_at)} and processing for verification
        </Typography>
      </CardContent>
    </CardX>
  )
}
