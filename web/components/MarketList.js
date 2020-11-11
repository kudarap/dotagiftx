import React from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import Avatar from '@material-ui/core/Avatar'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import green from '@material-ui/core/colors/lightGreen'
import { myMarket } from '@/service/api'
import { amount, dateFromNow } from '@/lib/format'
import Link from '@/components/Link'
import Button from '@/components/Button'
import BuyButton from '@/components/BuyButton'
import TableHeadCell from '@/components/TableHeadCell'
import ContactDialog from '@/components/ContactDialog'
import { MARKET_STATUS_REMOVED } from '@/constants/market'
import { retinaSrcSet } from '@/components/ItemImage'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'flex',
    padding: theme.spacing(2),
  },
  avatar: {
    marginRight: theme.spacing(1.5),
  },
  buyText: {
    color: green[600],
  },
}))

export default function MarketList({ data, error, currentUserID }) {
  const classes = useStyles()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))

  const [currentMarket, setCurrentMarket] = React.useState(null)
  const handleContactClick = marketIdx => {
    setCurrentMarket(data.data[marketIdx])
  }

  const router = useRouter()
  const handleRemoveClick = marketIdx => {
    const mktID = data.data[marketIdx].id
    ;(async () => {
      try {
        await myMarket.PATCH(mktID, { status: MARKET_STATUS_REMOVED })
        router.reload()
      } catch (e) {
        console.error(`Error: ${e.message}`)
      }
    })()
  }

  return (
    <>
      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableHeadCell>Seller</TableHeadCell>
              <TableHeadCell align="right">Price</TableHeadCell>
              {!isMobile && <TableHeadCell align="center" width={156} />}
            </TableRow>
          </TableHead>
          <TableBody>
            {error && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  Error retrieving data
                  <br />
                  <Typography variant="caption" color="textSecondary">
                    {error}
                  </Typography>
                </TableCell>
              </TableRow>
            )}

            {!error && !data && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  Loading...
                </TableCell>
              </TableRow>
            )}

            {!error && data.total_count === 0 && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  No offers available
                </TableCell>
              </TableRow>
            )}

            {data.data &&
              data.data.map((market, idx) => (
                <TableRow key={market.id} hover>
                  <TableCell component="th" scope="row" padding="none">
                    <Link
                      href="/profiles/[id]"
                      as={`/profiles/${market.user.steam_id}`}
                      disableUnderline>
                      <div className={classes.seller}>
                        <Avatar
                          className={classes.avatar}
                          alt={market.user.name}
                          {...retinaSrcSet(market.user.avatar, 40, 40)}
                        />
                        <div>
                          <strong>{market.user.name}</strong>
                          <br />
                          <Typography variant="caption" color="textSecondary">
                            {/* {market.user.steam_id} */}
                            Posted {dateFromNow(market.created_at)}
                          </Typography>
                        </div>
                      </div>
                    </Link>
                  </TableCell>
                  {!isMobile ? (
                    <>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {amount(market.price, market.currency)}
                        </Typography>
                      </TableCell>
                      <TableCell align="center">
                        {currentUserID === market.user.id ? (
                          // HOTFIX! wrapped button on div to prevent mixing up the styles(variant) of 2 buttons.
                          <div>
                            <Button variant="outlined" onClick={() => handleRemoveClick(idx)}>
                              Remove
                            </Button>
                          </div>
                        ) : (
                          <BuyButton variant="contained" onClick={() => handleContactClick(idx)}>
                            Contact Seller
                          </BuyButton>
                        )}
                      </TableCell>
                    </>
                  ) : (
                    <TableCell
                      align="right"
                      onClick={() =>
                        currentUserID === market.user.id
                          ? handleRemoveClick(idx)
                          : handleContactClick(idx)
                      }
                      style={{ cursor: 'pointer' }}>
                      <Typography variant="body2">${market.price.toFixed(2)}</Typography>
                      <Typography
                        variant="caption"
                        color="textSecondary"
                        style={{
                          color: currentUserID === market.user.id ? 'tomato' : '',
                        }}>
                        <u>{currentUserID === market.user.id ? 'Remove' : 'View'}</u>
                      </Typography>
                    </TableCell>
                  )}
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>
      <ContactDialog
        market={currentMarket}
        open={!!currentMarket}
        onClose={() => handleContactClick(null)}
      />
    </>
  )
}
MarketList.propTypes = {
  data: PropTypes.object.isRequired,
  error: PropTypes.string,
  currentUserID: PropTypes.string,
}
MarketList.defaultProps = {
  error: null,
  currentUserID: null,
}
