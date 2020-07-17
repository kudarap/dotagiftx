import React from 'react'
import PropTypes from 'prop-types'
import useSWR from 'swr'
import querystring from 'querystring'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import LinearProgress from '@material-ui/core/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import CatalogList from '@/components/CatalogList'
import TablePagination from '@/components/TablePagination'
import { CATALOGS, catalogSearch, fetcher } from '@/service/api'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(2.5),
  },
  listControl: {},
  paginator: {
    float: 'right',
  },
}))

const defaultFilter = {
  sort: 'created_at:desc',
  page: 1,
}

export default function Search({ catalogs: initialData }) {
  const classes = useStyles()

  const router = useRouter()
  const { query } = router
  query.page = Number(query.page || 1)
  const [filter, setFilter] = React.useState({
    ...defaultFilter,
    ...query,
  })

  const { data: items, error } = useSWR([CATALOGS, filter], fetcher, { initialData })
  React.useEffect(() => {
    setFilter({ ...filter, ...query })
  }, [query])

  const routerPush = f => {
    router.push(`/search?${querystring.stringify(f)}`)
  }

  const handlePageChange = (e, page) => {
    const f = { ...filter, page }
    setFilter(f)
    routerPush(f)
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          {filter.q && (
            <>
              <Typography component="h1" variant="h6">
                Results for &quot;{filter.q}&quot;
              </Typography>
              <br />
            </>
          )}

          {error && <div>failed to load</div>}
          {!items && <LinearProgress color="secondary" />}
          {!error && items && (
            <div>
              <CatalogList items={items.data} />
              <TablePagination
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
}

// This gets called on every request
export async function getServerSideProps({ query }) {
  const f = { ...defaultFilter, ...query }
  return {
    props: {
      catalogs: await catalogSearch(f),
    },
  }
}
