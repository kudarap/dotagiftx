import React from 'react'
import Head from 'next/head'
import Router, { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import LinearProgress from '@mui/material/LinearProgress'
import Select from '@mui/material/Select'
import type { SelectChangeEvent } from '@mui/material/Select'
import FormControl from '@mui/material/FormControl'
import MenuItem from '@mui/material/MenuItem'
import { catalogSearch, statsMarketSummaryOverall } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import CatalogList from '@/components/CatalogList'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import { APP_NAME, APP_URL } from '@/constants/strings'
import SearchInput from '@/components/SearchInput'
import * as format from '@/lib/format'
import type { SearchResponse, Item, MarketSummary } from '@/lib/types'

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

const sortOpts: Array<{ value: string; label: string }> = [
  ['popular', 'Most Popular'],
  ['recent', 'New Offers'],
  ['recent-bid', 'New Buy Orders'],
].map(([value, label]) => ({ value, label }))

function SelectSort({
  className = '',
  style = {},
  value,
  onChange,
}: {
  className?: string
  style?: React.CSSProperties
  value?: string
  onChange: (e: SelectChangeEvent<string>) => void
}) {
  return (
    <FormControl size="small" {...{ className, style }}>
      <Select id="select-sort" value={value} onChange={onChange}>
        {sortOpts.map(opt => (
          <MenuItem key={opt.value} value={opt.value}>
            {opt.label}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  )
}

interface SearchFilter {
  q?: string
  hero?: string
  origin?: string
  rarity?: string
  sort?: string
  page?: number
  [key: string]: string | number | undefined
}

export default function Search({
  catalogs: initialCatalogs,
  filter,
  canonicalURL,
}: {
  catalogs: SearchResponse<Item>
  filter: SearchFilter
  canonicalURL: string
}) {
  const { classes } = useStyles()

  const [catalogs, setCatalogs] = React.useState(initialCatalogs)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [sort, setSort] = React.useState(filter.sort)

  // Handle catalog request on page change.
  React.useEffect(() => {
    ;(async () => {
      setLoading(true)
      try {
        const res = (await catalogSearch(filter)) as SearchResponse<Item>
        setCatalogs(res)
      } catch (e) {
        setError((e as Error).message)
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
  const handleSelectSortChange = (e: SelectChangeEvent<string>) => {
    const nextSort = e.target.value as string
    setSort(nextSort)
    router.push({
      pathname: '/search',
      query: {
        ...filter,
        sort: nextSort,
      },
    })
  }

  const isBidType = filter.sort === 'recent-bid'

  const handleSearchSubmit = (keyword: string) => {
    Router.push(`/search?q=${keyword}`)
  }

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
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
                  {catalogs && catalogs.total_count}&nbsp;results for &quot;{searchTerm}&quot;
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
                error={error ?? undefined}
                bidType={isBidType}
              />
              {!error && (
                <TablePaginationRouter
                  linkProps={linkProps}
                  colSpan={3}
                  style={{ textAlign: 'right' }}
                  count={catalogs.total_count}
                  page={filter.page as number}
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

const catalogSearchFilter: SearchFilter = { sort: 'popular', page: 1 }

// This gets called on every request
export async function getServerSideProps({ query }: { query: Record<string, string | undefined> }) {
  const filter: SearchFilter = { ...catalogSearchFilter, ...query }
  filter.page = Number(query.page || 1)

  let catalogs = {} as SearchResponse<Item>
  let error: string | null = null
  try {
    catalogs = (await catalogSearch(filter)) as SearchResponse<Item>
  } catch (e) {
    error = (e as Error).message
  }

  const marketSummary = (await statsMarketSummaryOverall()) as MarketSummary
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
