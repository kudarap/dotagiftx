import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import InputBase from '@material-ui/core/InputBase'
import InputAdornment from '@material-ui/core/InputAdornment'
import LinearProgress from '@material-ui/core/LinearProgress'
import Paper from '@material-ui/core/Paper'
import SearchIcon from '@material-ui/icons/Search'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { myMarketSearch } from '@/service/api'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import MyMarketList from '@/components/MyMarketList'
import TablePagination from '@/components/TablePagination'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(1),
    },
    marginTop: theme.spacing(4),
  },
  searchPaper: {
    [theme.breakpoints.down('xs')]: {
      width: '100%',
    },
    float: 'right',
    flexGrow: 1,
    padding: '3px 12px',
    marginBottom: 2,
    display: 'flex',
    alignItems: 'center',
    fontStyle: 'italic',
    // '&:hover': {
    //   borderColor: theme.palette.grey[700],
    // },
    backgroundColor: theme.palette.primary.dark,
    opacity: 0.8,
  },
  searchInput: {
    marginLeft: theme.spacing(1),
  },
  searchIcon: {
    color: theme.palette.grey[500],
  },
}))

const activeMarketFilter = {
  status: MARKET_STATUS_LIVE,
  sort: 'created_at:desc',
  page: 1,
}

const initialDatatable = {
  data: [],
  result_count: 0,
  total_count: 0,
  loading: true,
  error: null,
}

export default function MyListings() {
  const classes = useStyles()

  const [activeLists, setActiveLists] = React.useState(initialDatatable)
  const [filter, setFilter] = React.useState(activeMarketFilter)

  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await myMarketSearch(filter)
        setActiveLists({ ...activeLists, loading: false, ...res })
      } catch (e) {
        setActiveLists({ ...activeLists, loading: false, error: e.message })
      }
    })()
  }, [filter])

  const handlePageChange = (e, page) => {
    setFilter({ ...filter, page })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <div>
            <Typography component="span" gutterBottom>
              My Active Listings
            </Typography>
            <Paper className={classes.searchPaper} elevation={0}>
              <SearchIcon className={classes.searchIcon} />
              <InputBase
                fullWidth
                className={classes.searchInput}
                color="secondary"
                placeholder="Search Listings"
                style={{ float: 'right' }}
              />
            </Paper>
          </div>

          {activeLists.error && <div>failed to load active listings</div>}
          {activeLists.loading && <LinearProgress color="secondary" />}
          <MyMarketList datatable={activeLists} />
          <TablePagination
            style={{ textAlign: 'right' }}
            count={activeLists.total_count || 0}
            page={filter.page}
            onChangePage={handlePageChange}
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
