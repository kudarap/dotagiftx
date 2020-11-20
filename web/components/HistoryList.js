import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import { makeStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import * as format from '@/lib/format'
import Button from '@/components/Button'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'
import HistoryViewDialog from '@/components/HistoryViewDialog'
import AppContext from '@/components/AppContext'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'inline-flex',
  },
  item: {
    padding: theme.spacing(2, 2, 2, 0),
    display: 'flex',
    cursor: 'pointer',
  },
  image: {
    margin: theme.spacing(-1, 1, -1, 1),
    width: 77,
    height: 55,
  },
}))

export default function HistoryList({ datatable, error }) {
  const classes = useStyles()
  const { isMobile } = useContext(AppContext)

  const [currentMarket, setCurrentMarket] = React.useState(null)

  const handleUpdateClick = marketIdx => {
    setCurrentMarket(datatable.data[marketIdx])
  }

  if (error) {
    return <p>Error</p>
  }

  if (!datatable) {
    return <p>Loading...</p>
  }

  return (
    <>
      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableHeadCell>
                Items ({format.numberWithCommas(datatable.total_count)})
              </TableHeadCell>
              {!isMobile && (
                <>
                  <TableHeadCell align="right">Updated</TableHeadCell>
                  <TableHeadCell align="right">Price</TableHeadCell>
                </>
              )}
              <TableHeadCell align="center" width={70} />
            </TableRow>
          </TableHead>
          <TableBody>
            {datatable.data &&
              datatable.data.map((market, idx) => (
                <TableRow key={market.id} hover>
                  <TableCell
                    component="th"
                    scope="row"
                    padding="none"
                    className={classes.item}
                    onClick={() => handleUpdateClick(idx)}>
                    <ItemImage
                      className={classes.image}
                      image={market.item.image}
                      width={77}
                      height={55}
                      title={market.item.name}
                      rarity={market.item.rarity}
                    />
                    <div>
                      <strong>{market.item.name}</strong>
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {market.item.hero}
                      </Typography>
                      <RarityTag rarity={market.item.rarity} />
                    </div>
                  </TableCell>
                  {!isMobile ? (
                    <>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {moment(market.updated_at).fromNow()}
                        </Typography>
                      </TableCell>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {format.amount(market.price, market.currency)}
                        </Typography>
                      </TableCell>
                      <TableCell align="center">
                        <Button variant="outlined" onClick={() => handleUpdateClick(idx)}>
                          View
                        </Button>
                      </TableCell>
                    </>
                  ) : (
                    <TableCell
                      align="right"
                      onClick={() => handleUpdateClick(idx)}
                      style={{ cursor: 'pointer' }}>
                      <Typography variant="body2" color="secondary">
                        {format.amount(market.price, market.currency)}
                      </Typography>
                      <Typography variant="caption" color="textSecondary" noWrap>
                        {moment(market.updated_at).fromNow()}
                      </Typography>
                    </TableCell>
                  )}
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>
      <HistoryViewDialog
        open={!!currentMarket}
        market={currentMarket}
        onClose={() => handleUpdateClick(null)}
      />
    </>
  )
}
HistoryList.propTypes = {
  datatable: PropTypes.object.isRequired,
  error: PropTypes.string,
}
HistoryList.defaultProps = {
  error: null,
}
