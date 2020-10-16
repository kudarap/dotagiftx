import React from 'react'
import PropTypes from 'prop-types'
import debounce from 'lodash/debounce'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import InputBase from '@material-ui/core/InputBase'
import LinearProgress from '@material-ui/core/LinearProgress'
import SearchIcon from '@material-ui/icons/Search'
import Typography from '@material-ui/core/Typography'
import * as format from '@/lib/format'
import Button from '@/components/Button'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'
import ItemImage from '@/components/ItemImage'
import MarketUpdateDialog from '@/components/MarketUpdateDialog'

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
  searchPaper: {
    padding: '3px 12px',
    display: 'flex',
    alignItems: 'center',
    backgroundColor: theme.palette.primary.dark,
    opacity: 0.8,
    margin: theme.spacing(1),
  },
  searchInput: {
    fontSize: theme.typography.body2.fontSize,
    marginLeft: theme.spacing(1),
  },
  searchIcon: {
    color: theme.palette.grey[500],
  },
}))

export default function MyMarketList({ datatable, loading, error, onSearchInput }) {
  const classes = useStyles()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))

  const [currentMarket, setCurrentMarket] = React.useState(null)

  const handleUpdateClick = marketIdx => {
    setCurrentMarket(datatable.data[marketIdx])
  }

  const debounceSearch = debounce(onSearchInput, 500)
  const handleSearchInput = e => {
    debounceSearch(e.target.value)
  }

  return (
    <>
      <TableContainer component={Paper}>
        {loading && <LinearProgress color="secondary" />}
        <Table className={classes.table} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableHeadCell padding="none" colSpan={isMobile ? 2 : 1}>
                <Paper className={classes.searchPaper} elevation={0}>
                  <SearchIcon className={classes.searchIcon} />
                  <InputBase
                    fullWidth
                    onInput={handleSearchInput}
                    className={classes.searchInput}
                    color="secondary"
                    placeholder="Filter items"
                  />
                </Paper>
              </TableHeadCell>
              {!isMobile && (
                <>
                  <TableHeadCell align="right">Listed</TableHeadCell>
                  <TableHeadCell align="right">Price</TableHeadCell>
                  <TableHeadCell align="center" width={70} />
                </>
              )}
            </TableRow>
          </TableHead>
          <TableBody>
            {error && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  Error retrieving data
                  <br />
                  <Typography variant="caption" color="textSecondary">
                    {format.errorSimple(error)}
                  </Typography>
                </TableCell>
              </TableRow>
            )}

            {!error && datatable.data.length === 0 && (
              <TableRow>
                <TableCell align="center" colSpan={3}>
                  No Result
                </TableCell>
              </TableRow>
            )}

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
                      image={`/200x100/${market.item.image}`}
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
                  {!isMobile && (
                    <>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {format.dateFromNow(market.created_at)}
                        </Typography>
                      </TableCell>
                      <TableCell align="right">
                        <Typography variant="body2">
                          {format.amount(market.price, market.currency)}
                        </Typography>
                      </TableCell>
                    </>
                  )}
                  <TableCell align="center">
                    <Button variant="outlined" onClick={() => handleUpdateClick(idx)}>
                      Update
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </TableContainer>
      <MarketUpdateDialog
        open={!!currentMarket}
        market={currentMarket}
        onClose={() => handleUpdateClick(null)}
      />
    </>
  )
}
MyMarketList.propTypes = {
  datatable: PropTypes.object.isRequired,
  onSearchInput: PropTypes.func,
  loading: PropTypes.bool,
  error: PropTypes.string,
}
MyMarketList.defaultProps = {
  onSearchInput: () => {},
  loading: false,
  error: null,
}
