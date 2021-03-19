import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Card from '@material-ui/core/Card'
import CardActions from '@material-ui/core/CardActions'
import CardContent from '@material-ui/core/CardContent'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import { LightTheme } from '@/components/Theme'
import {
  VERIFIED_DELIVERY_MAP_LABEL,
  VERIFIED_DELIVERY_MAP_TEXT,
  VERIFIED_INVENTORY_MAP_LABEL,
  VERIFIED_INVENTORY_MAP_TEXT,
} from '@/constants/verified'
import { dateFromNow } from '@/lib/format'
import Link from '@/components/Link'
import { Popover } from '@material-ui/core'

const useStyles = makeStyles(theme => ({
  root: {
    minWidth: 300,
  },
  link: {},
  image: {},
}))

const assetModifier = asset => {
  let isGiftable = asset.gift_once ? 'Yes' : 'No'
  if (asset.type.startsWith('Immortal')) {
    isGiftable = '?'
  }

  return { ...asset, isGiftable }
}

const getInventoryURL = steamID => `https://steamcommunity.com/profiles/${steamID}/inventory/#570_2`

export default function VerifiedStatusCard({ market, ...other }) {
  const classes = useStyles()

  if (market === null) {
    return null
  }

  const { inventory, delivery } = market
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
    return null
  }

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
                  {!isDelivery && <TableCell align="center">Giftable</TableCell>}
                  {isDelivery && <TableCell>From</TableCell>}
                  {isDelivery && <TableCell>Received</TableCell>}
                </TableRow>
              </TableHead>
              <TableBody>
                {source.steam_assets.map(assetModifier).map(asset => (
                  <TableRow key={asset.name}>
                    <TableCell component="th" scope="row">
                      <Link
                        target="_blank"
                        rel="noreferrer noopener"
                        href={`${inventoryURL}_${asset.asset_id}`}>
                        {asset.name}
                      </Link>
                    </TableCell>
                    {!isDelivery && <TableCell align="center">{asset.isGiftable}</TableCell>}
                    {isDelivery && <TableCell>{asset.gift_from}</TableCell>}
                    {isDelivery && <TableCell>{asset.date_received}</TableCell>}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </>
        )}
      </CardContent>
      <CardActions style={{ float: 'right' }}>
        <Link
          variant="caption"
          color="primary"
          target="_blank"
          rel="noreferrer noopener"
          href="https://steaminventory.org">
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
    <LightTheme>
      <Card {...props} />
    </LightTheme>
  )
}

export function VerifiedStatusPopover({ market, ...other }) {
  return (
    <Popover
      style={{ marginLeft: 4 }}
      anchorOrigin={{
        vertical: 'top',
        horizontal: 'left',
      }}
      transformOrigin={{
        vertical: 'top',
        horizontal: 'left',
      }}
      {...other}>
      <VerifiedStatusCard market={market} onMouseLeave={other.onClose} />
    </Popover>
  )
}
VerifiedStatusPopover.propTypes = VerifiedStatusCard.propTypes
VerifiedStatusPopover.defaultProps = VerifiedStatusCard.defaultProps
