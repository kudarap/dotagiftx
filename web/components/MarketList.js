import React from 'react'
import PropTypes from 'prop-types'
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
import { CDN_URL } from '@/service/api'
import Link from '@/components/Link'
import Button from '@/components/Button'
import BuyButton from '@/components/BuyButton'
import TableHeadCell from '@/components/TableHeadCell'
import ContactDialog from '@/components/ContactDialog'
import { amount } from '@/lib/format'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'flex',
    padding: theme.spacing(2),
  },
  avatar: {
    marginRight: theme.spacing(1.5),
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

  return (
    <>
      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableHeadCell>Seller</TableHeadCell>
              <TableHeadCell align="right">Price</TableHeadCell>
              <TableHeadCell align="center" width={156} />
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

            {data.data &&
              data.data.map((market, idx) => (
                <TableRow key={market.id} hover>
                  <TableCell component="th" scope="row" padding="none">
                    <Link href="/user/[id]" as={`/user/${market.user.steam_id}`} disableUnderline>
                      <div className={classes.seller}>
                        {!isMobile && (
                          <Avatar
                            className={classes.avatar}
                            src={`${CDN_URL}/${market.user.avatar}`}
                          />
                        )}
                        <div>
                          <strong>{market.user.name}</strong>
                          <br />
                          <Typography variant="caption" color="textSecondary">
                            {market.user.steam_id}
                          </Typography>
                        </div>
                      </div>
                    </Link>
                  </TableCell>
                  <TableCell align="right">
                    <Typography variant="body2">{amount(market.price, market.currency)}</Typography>
                  </TableCell>
                  <TableCell align="center">
                    {currentUserID === market.user.id ? (
                      <Button variant="outlined">Remove</Button>
                    ) : (
                      <BuyButton variant="contained" onClick={() => handleContactClick(idx)}>
                        Contact Seller
                      </BuyButton>
                    )}
                  </TableCell>
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
