import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Select from '@material-ui/core/Select'
import FormControl from '@material-ui/core/FormControl'
import MenuItem from '@material-ui/core/MenuItem'
import { catalogSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import CatalogList from '@/components/CatalogList'
import TablePaginationRouter from '@/components/TablePaginationRouter'
import { APP_NAME, APP_URL } from '@/constants/strings'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(2),
  },
  listControl: {},
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

export default function Search({ catalogs: initialCatalogs, filter, canonicalURL }) {
  const classes = useStyles()

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

  return (
    <>
      <Head>
        <title>{metaTitle}</title>
        <meta name="description" content={metaDesc} />
        <link rel="canonical" href={canonicalURL} />
      </Head>
      <Header />

      <main className={classes.main}>
        <Container>
          {searchTerm && (
            <>
              <Typography component="h1" variant="h6">
                {catalogs && catalogs.total_count} results for &quot;{searchTerm}&quot;
              </Typography>
            </>
          )}

          <SelectSort style={{ float: 'right' }} value={sort} onChange={handleSelectSortChange} />
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

  const canonicalURL = `${APP_URL}/search?q=${filter.q}`

  return {
    props: {
      canonicalURL,
      filter,
      catalogs,
      error,
    },
  }
}
