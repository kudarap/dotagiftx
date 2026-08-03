import React from 'react'
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
import type { PopperProps } from '@mui/material/Popper'
import { indigo } from '@mui/material/colors'
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
import type { Market, InventoryAsset } from '@/lib/types'

const useStyles = makeStyles()({
  root: {
    minWidth: 300,
    zIndex: 1,
    marginTop: 18,
  },
  poweredBy: {
    color: indigo[400],
  },
})

type AssetView = InventoryAsset & {
  isGiftable: string
  displayName: string | undefined
  received: string | undefined
}

const assetModifier = (asset: InventoryAsset): AssetView => {
  let isGiftable = asset.gift_once ? 'Yes' : 'No'
  if (asset.type && asset.type.startsWith('Immortal') && !asset.gift_once) {
    isGiftable = '?'
  }

  // Shows gift containing items.
  let displayName: string | undefined = asset.name
  let received: string | undefined = asset.date_received
  if (asset.name === 'Wrapped Gift') {
    displayName = asset.contains
    received = asset.name
  }

  return { ...asset, isGiftable, displayName, received }
}

const getInventoryURL = (steamID: string) =>
  `https://steamcommunity.com/profiles/${steamID}/inventory/#570_2`

const formatDuration = (ms: number) => {
  let millis = ms
  if (millis < 0) {
    millis = -millis
  }
  const time = {
    minute: Math.floor(millis / 60000) % 60,
    second: Math.floor(millis / 1000) % 60,
    m: Math.floor(millis) % 1000,
  }

  return Object.entries(time)
    .filter(val => val[1] !== 0)
    .map(([key, val]) => `${val} ${key}${val !== 1 ? 's' : ''}`)
    .join(', ')
}

function CardX(props: React.ComponentProps<typeof Card>) {
  return (
    <Card
      // dark carnival event theme
      sx={{ mt: 2.25, boxShadow: 10, background: '#292e3d', border: '1px solid #515051' }}
      {...props}
    />
  )
}

interface VerifiedStatusCardProps {
  market: Market | null
  onMouseLeave?: () => void
}

export default function VerifiedStatusCard({
  market,
  onMouseLeave,
  ...other
}: VerifiedStatusCardProps) {
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
    inventoryURL = getInventoryURL(market.partner_steam_id || '')
  }
  if (!source) {
    if (market.resell) {
      return <ResellCard data={market} />
    }

    return null
  }

  const steamInvProfile =
    market.status === MARKET_STATUS_SOLD ? market.partner_steam_id || '' : market.user.steam_id

  return (
    <CardX className={classes.root} onMouseLeave={onMouseLeave} {...other}>
      <CardContent>
        <Typography variant="h5" component="h2">
          {mapLabel[source.status]}
        </Typography>
        <Typography color="textSecondary" variant="caption" component="p" sx={{ mb: 1 }}>
          Processed {dateFromNow(source.updated_at)}
          {source?.elapsed_ms ? <span>&nbsp;in {formatDuration(source.elapsed_ms)}</span> : null}
        </Typography>

        <Typography component="p">{mapText[source.status]}.</Typography>

        {source.steam_assets && (
          <>
            {!isDelivery && (
              <Typography variant="body2">
                Found <strong>{source.bundle_count}</strong> bundle
                {source.bundle_count && source.bundle_count > 1 && 's'}
              </Typography>
            )}
            <Table size="small">
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
                        href={`${inventoryURL}_${asset.asset_id}`}
                      >
                        <strong>{asset.displayName || asset.name}</strong>
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
        {source?.verified_by ? (
          <Typography variant="caption">
            Powered by <strong style={{ textTransform: 'capitalize' }}>{source.verified_by}</strong>
          </Typography>
        ) : (
          <Link
            className={classes.poweredBy}
            variant="caption"
            target="_blank"
            rel="noreferrer noopener"
            underline="none"
            href={`https://steaminventory.org/?profile=${steamInvProfile}`}
          >
            Powered by <strong>SteamInventory.org</strong>
          </Link>
        )}
      </CardActions>
    </CardX>
  )
}

export function VerifiedStatusPopover({
  market,
  onClose,
  ...other
}: { market: Market | null; onClose?: () => void } & Omit<PopperProps, 'onClose'>) {
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
      {...other}
    >
      <VerifiedStatusCard market={market} onMouseLeave={onClose} />
    </Popper>
  )
}

function ClearedGift() {
  return (
    <Typography color="textSecondary" variant="body2" component="em">
      --
      {/* cleared */}
    </Typography>
  )
}

function ResellCard({ data }: { data: Market }) {
  return (
    <CardX>
      <CardContent>
        <Typography variant="h5" component="h2">
          Item Resell
        </Typography>
        <Typography color="textSecondary" variant="caption" component="p" sx={{ mb: 1 }}>
          Verified by <strong>{data.user.name}</strong> on {dateFromNow(data.created_at || '')}.
        </Typography>
        <Typography component="p">
          Item manually verified by seller from partner&apos;s inventory
        </Typography>
      </CardContent>
    </CardX>
  )
}

function PendingCard({ data }: { data: Market }) {
  return (
    <CardX>
      <CardContent>
        <Typography variant="h5" component="h2">
          Pending
        </Typography>
        <Typography color="textSecondary" variant="caption" component="p" sx={{ mb: 1 }}>
          Posted {dateFromNow(data.created_at || '')} and processing for verification.
        </Typography>
      </CardContent>
    </CardX>
  )
}
