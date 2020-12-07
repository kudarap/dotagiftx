import React, { useContext, useEffect } from 'react'
import PropTypes from 'prop-types'
import has from 'lodash/has'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import { myMarket } from '@/service/api'
import { amount, dateFromNow } from '@/lib/format'
import Link from '@/components/Link'
import Button from '@/components/Button'
import BuyButton from '@/components/BuyButton'
import TableHeadCell from '@/components/TableHeadCell'
import ContactDialog from '@/components/ContactDialog'
import ContactBuyerDialog from '@/components/ContactBuyerDialog'
import { MARKET_STATUS_REMOVED } from '@/constants/market'
import { retinaSrcSet } from '@/components/ItemImage'
import AppContext from '@/components/AppContext'
import { Tab, Tabs } from '@material-ui/core'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'flex',
    padding: theme.spacing(2),
  },
  avatar: {
    marginRight: theme.spacing(1.5),
  },
  tableHead: {
    // background: theme.palette.grey[900],
    background: '#202a2f',
  },
  tabs: {
    '& .MuiTabs-indicator': {
      background: theme.palette.grey[100],
    },
  },
}))

export default function MarketList({ offers, buyOrders, error, loading, pagination }) {
  const classes = useStyles()
  const { isMobile, currentAuth } = useContext(AppContext)
  const currentUserID = currentAuth.user_id || null

  const [tabIdx, setTabIdx] = React.useState(0)

  const router = useRouter()
  useEffect(() => {
    if (has(router.query, 'buyorder')) {
      setTabIdx(1)
    } else {
      setTabIdx(0)
    }
  }, [router.query])

  const handleTabChange = (e, value) => {
    setTabIdx(value)
    let p = `/${router.query.slug}`
    if (value === 1) {
      p += '?buyorder'
    }

    router.push(p)
  }

  const [currentMarket, setCurrentMarket] = React.useState(null)
  const handleContactClick = marketIdx => {
    let src = offers
    if (tabIdx === 1) {
      src = buyOrders
    }

    setCurrentMarket(src.data[marketIdx])
  }

  const handleRemoveClick = marketIdx => {
    let src = offers
    if (tabIdx === 1) {
      src = buyOrders
    }

    const mktID = src.data[marketIdx].id
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
        <Table className={classes.table} aria-label="market list table">
          <TableHead className={classes.tableHead}>
            <TableRow>
              <TableHeadCell colSpan={3} padding="none">
                <Tabs
                  className={classes.tabs}
                  variant="fullWidth"
                  value={tabIdx}
                  onChange={handleTabChange}>
                  <Tab
                    value={0}
                    label={`${offers.total_count || ''} Offers`}
                    style={{ textTransform: 'none' }}
                  />
                  <Tab
                    value={1}
                    label={`${buyOrders.total_count || ''} Buy Orders`}
                    style={{ textTransform: 'none' }}
                  />
                </Tabs>
              </TableHeadCell>
            </TableRow>
          </TableHead>

          {tabIdx === 0 ? (
            <TableBody style={{ opacity: loading ? 0.5 : 1 }}>
              <TableRow>
                <TableHeadCell size="small">
                  <Typography color="textSecondary" variant="body2">
                    Seller
                  </Typography>
                </TableHeadCell>
                <TableHeadCell size="small" align="right">
                  <Typography color="textSecondary" variant="body2">
                    Price
                  </Typography>
                </TableHeadCell>
                {!isMobile && <TableHeadCell size="small" align="center" width={160} />}
              </TableRow>

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

              {!offers.data && (
                <TableRow>
                  <TableCell align="center" colSpan={3}>
                    Loading...
                  </TableCell>
                </TableRow>
              )}

              {!error && offers.total_count === 0 && (
                <TableRow>
                  <TableCell align="center" colSpan={3}>
                    No offers available
                  </TableCell>
                </TableRow>
              )}

              {offers.data &&
                offers.data.map((market, idx) => (
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
          ) : (
            <TableBody style={{ opacity: loading ? 0.5 : 1 }}>
              <TableRow>
                <TableHeadCell size="small">
                  <Typography color="textSecondary" variant="body2">
                    Buyer
                  </Typography>
                </TableHeadCell>
                <TableHeadCell size="small" align="right">
                  <Typography color="textSecondary" variant="body2">
                    Price
                  </Typography>
                </TableHeadCell>
                {!isMobile && <TableHeadCell size="small" align="center" width={160} />}
              </TableRow>

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

              {!buyOrders.data && (
                <TableRow>
                  <TableCell align="center" colSpan={3}>
                    Loading...
                  </TableCell>
                </TableRow>
              )}

              {!error && buyOrders.total_count === 0 && (
                <TableRow>
                  <TableCell align="center" colSpan={3}>
                    No buy orders
                  </TableCell>
                </TableRow>
              )}

              {buyOrders.data &&
                buyOrders.data.map((market, idx) => (
                  <TableRow key={market.id} hover>
                    <TableCell component="th" scope="row" padding="none">
                      {/* Check for redacted display */}
                      {market.user.id ? (
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
                      ) : (
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
                      )}
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
                            <BuyButton
                              // Check for redacted user and disabled them for opening the dialog.
                              disabled={!market.user.id}
                              color="primary"
                              variant="contained"
                              onClick={() => handleContactClick(idx)}>
                              {market.user.id ? (
                                'Contact Buyer'
                              ) : (
                                <Typography variant="caption" color="textSecondary">
                                  Sign in to view
                                </Typography>
                              )}
                            </BuyButton>
                          )}
                        </TableCell>
                      </>
                    ) : (
                      <TableCell
                        align="right"
                        onClick={() => {
                          if (!market.user.id) {
                            return
                          }

                          if (currentUserID === market.user.id) {
                            handleRemoveClick(idx)
                            return
                          }

                          handleContactClick(idx)
                        }}
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
          )}
        </Table>
      </TableContainer>

      {/* Only display pagination on offers */}
      {tabIdx === 0 && pagination}

      {tabIdx === 1 && buyOrders.total_count > 10 && (
        <Typography color="textSecondary" align="right" variant="body2" style={{ margin: 8 }}>
          `${buyOrders.total_count - 10} more hidden buy orders at $
          {amount(buyOrders.data[9].price, 'USD')} or less`
        </Typography>
      )}

      {((tabIdx === 0 && !pagination) || (tabIdx === 1 && buyOrders.total_count <= 10)) && (
        <div style={{ margin: 8 }}>&nbsp;</div>
      )}

      <ContactDialog
        market={currentMarket}
        open={tabIdx === 0 && !!currentMarket}
        onClose={() => handleContactClick(null)}
      />

      <ContactBuyerDialog
        market={currentMarket}
        open={tabIdx === 1 && !!currentMarket}
        onClose={() => handleContactClick(null)}
      />
    </>
  )
}
MarketList.propTypes = {
  offers: PropTypes.object.isRequired,
  buyOrders: PropTypes.object.isRequired,
  pagination: PropTypes.element,
  error: PropTypes.string,
}
MarketList.defaultProps = {
  pagination: null,
  error: null,
}
