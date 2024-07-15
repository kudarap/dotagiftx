import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import Router, { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import LinearProgress from '@mui/material/LinearProgress'
import Select from '@mui/material/Select'
import FormControl from '@mui/material/FormControl'
import MenuItem from '@mui/material/MenuItem'
import { catalogSearch, statsMarketSummary } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import CatalogList from '@/components/CatalogList'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import { APP_NAME, APP_URL } from '@/constants/strings'
import SearchInput from '@/components/SearchInput'
import * as format from '@/lib/format'

const useStyles = makeStyles()(theme => ({
  main: {
    marginTop: theme.spacing(2),
  },
  searchBar: {
    display: 'flex',
    justifyContent: 'space-between',
    padding: theme.spacing(1, 0, 1),
  },
  paginator: {
    float: 'right',
  },
}))

const sortOpts = [
  ['popular', 'Most Popular'],
  ['recent', 'New Offers'],
  ['recent-bid', 'New Buy Orders'],
].map(([value, label]) => ({ value, label }))

function SelectSort({ className, style, ...other }) {
  return (
    <FormControl size="small" {...{ className, style }}>
      <Select id="select-sort" {...other}>
        {sortOpts.map(opt => (
          <MenuItem value={opt.value}>{opt.label}</MenuItem>
        ))}
      </Select>
    </FormControl>
  )
}

export default function Search({ catalogs: initialCatalogs, marketSummary, filter, canonicalURL }) {
  const { classes } = useStyles()

  const [catalogs, setCatalogs] = React.useState(initialCatalogs)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState(null)
  const [sort, setSort] = React.useState(filter.sort)

  // Handle catalog request on page change.
  React.useEffect(() => {
    ;(async () => {
      setLoading(true)
      try {
        const res = await catalogSearch(filter)
        setCatalogs(res)
      } catch (e) {
        setError(e.message)
      }
      setLoading(false)
    })()
  }, [filter])

  let metaTitle = `${APP_NAME} :: Search`
  let metaDesc = `Search for item name, hero, treasure`
  const searchTerm = filter.q || filter.hero || filter.origin || filter.rarity
  if (searchTerm) {
    metaTitle += ` ${searchTerm}`
    metaDesc = `${catalogs && catalogs.total_count} results for "${searchTerm}"`
  }

  const linkProps = { href: '/search', query: filter }

  const router = useRouter()
  const handleSelectSortChange = e => {
    setSort(e.target.value)
    linkProps.query.sort = e.target.value
    router.push(linkProps)
  }

  const isBidType = filter.sort === 'recent-bid'

  const handleSearchSubmit = keyword => {
    Router.push(`/search?q=${keyword}`)
  }

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={canonicalURL} />
      </Head>
      <Header />

      <main className={classes.main}>
        <Container>
          <SearchInput value={searchTerm} onSubmit={handleSearchSubmit} label="" />

          <div className={classes.searchBar}>
            {searchTerm && (
              <div>
                <Typography component="h1" variant="h6">
                  {catalogs && catalogs.total_count} results for &quot;{searchTerm}&quot;
                </Typography>
              </div>
            )}

            <SelectSort style={{ float: 'right' }} value={sort} onChange={handleSelectSortChange} />
          </div>

          {!catalogs && <LinearProgress color="secondary" />}
          {catalogs && (
            <div>
              <CatalogList
                items={catalogs.data}
                loading={loading}
                error={error}
                bidType={isBidType}
              />
              {!error && (
                <TablePaginationRouter
                  linkProps={linkProps}
                  colSpan={3}
                  style={{ textAlign: 'right' }}
                  count={catalogs.total_count}
                  page={filter.page}
                />
              )}
            </div>
          )}
        </Container>
      </main>

      <Footer />
    </>
  )
}
Search.propTypes = {
  catalogs: PropTypes.object.isRequired,
  filter: PropTypes.object.isRequired,
  canonicalURL: PropTypes.string.isRequired,
}

const catalogSearchFilter = { sort: 'popular', page: 1 }

// This gets called on every request
export async function getServerSideProps({ query }) {
  const filter = { ...catalogSearchFilter, ...query }
  filter.page = Number(query.page || 1)

  let catalogs = {}
  let error = null
  try {
    catalogs = await catalogSearch(filter)
  } catch (e) {
    error = e.message
  }

  const marketSummary = await statsMarketSummary()
  marketSummary.live = format.numberWithCommas(marketSummary.live)
  marketSummary.reserved = format.numberWithCommas(marketSummary.reserved)
  marketSummary.sold = format.numberWithCommas(marketSummary.sold)
  marketSummary.bids.bid_live = format.numberWithCommas(marketSummary.bids.bid_live)

  const canonicalURL = `${APP_URL}/search?q=${filter.q}`

  return {
    props: {
      canonicalURL,
      filter,
      catalogs,
      marketSummary,
      error,
    },
  }
}
