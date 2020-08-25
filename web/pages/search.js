import React from 'react'
import PropTypes from 'prop-types'
import has from 'lodash/has'
import isEqual from 'lodash/isEqual'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import { catalogSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import CatalogList from '@/components/CatalogList'
import TablePaginationRouter from '@/components/TablePaginationRouter'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(2.5),
  },
  listControl: {},
  paginator: {
    float: 'right',
  },
}))

export default function Search({ catalogs: initialCatalogs, filter: initialFilter }) {
  const classes = useStyles()

  const [filter, setFilter] = React.useState(initialFilter)
  const [page, setPage] = React.useState(filter.page)
  const [catalogs, setCatalogs] = React.useState(initialCatalogs)
  const [error, setError] = React.useState(null)

  // Handle search keyword change and resets page if available.
  React.useEffect(() => {
    if (isEqual(filter, initialFilter)) {
      return
    }

    if (has(initialFilter, 'q')) {
      setFilter({ ...filter, q: initialFilter.q, page: 1 })
      setPage(1)
    }
  }, [initialFilter])

  // Handle catalog request on page change.
  React.useEffect(() => {
    ;(async () => {
      try {
        const res = await catalogSearch(filter)
        setCatalogs(res)
      } catch (e) {
        setError(e.message)
      }
    })()
  }, [filter])

  const handlePageChange = (e, p) => {
    setFilter({ ...filter, page: p })
    setPage(p)
  }

  const linkProps = { href: '/search', query: filter }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          {filter.q && (
            <>
              <Typography component="h1" variant="h6">
                {catalogs && catalogs.total_count} results for &quot;{filter.q}&quot;
              </Typography>
              <br />
            </>
          )}

          {!catalogs && <LinearProgress color="secondary" />}
          {catalogs && (
            <div>
              <CatalogList items={catalogs.data} error={error} />
              <TablePaginationRouter
                linkProps={linkProps}
                colSpan={3}
                style={{ textAlign: 'right' }}
                count={catalogs.total_count}
                page={page}
                onChangePage={handlePageChange}
              />
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
}

const catalogSearchFilter = { sort: 'created_at:desc', page: 1 }

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

  return {
    props: {
      filter,
      catalogs,
      error,
    },
  }
}
