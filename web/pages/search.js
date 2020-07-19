import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import { catalogSearch } from '@/service/api'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import CatalogList from '@/components/CatalogList'
import TablePagination from '@/components/TablePaginationRouter'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(2.5),
  },
  listControl: {},
  paginator: {
    float: 'right',
  },
}))

const defaultFilter = { sort: 'created_at:desc', page: 1 }

export default function Search({ catalogs: items, filter: propFilter }) {
  const classes = useStyles()

  const [filter, setFilter] = React.useState(propFilter)

  React.useEffect(() => {
    setFilter({ ...filter, ...propFilter })
  }, [propFilter])

  const handlePageChange = (e, page) => {
    const f = { ...filter, page }
    setFilter(f)
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
                {items && items.total_count} results for &quot;{filter.q}&quot;
              </Typography>
              <br />
            </>
          )}

          {!items && <LinearProgress color="secondary" />}
          {items && (
            <div>
              <CatalogList items={items.data} />
              <TablePagination
                linkProps={linkProps}
                colSpan={3}
                style={{ textAlign: 'right' }}
                page={filter.page}
                count={items.total_count}
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

// This gets called on every request
export async function getServerSideProps({ query }) {
  const filter = { ...defaultFilter, ...query }
  filter.page = Number(query.page || 1)

  return {
    props: {
      catalogs: await catalogSearch(filter),
      filter,
    },
  }
}
