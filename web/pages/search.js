import React from 'react'
import useSWR from 'swr'
import querystring from 'querystring'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import LinearProgress from '@material-ui/core/LinearProgress'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import ItemList from '@/components/ItemList'
import SearchInput from '@/components/SearchInput'
import TablePagination from '@/components/TablePagination'
import { CATALOGS, fetcher } from '@/service/api'

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
  sort: 'name',
  page: 1,
}

export default function Search() {
  const classes = useStyles()

  const router = useRouter()
  const { query } = router
  query.page = Number(query.page || 1)
  const [filter, setFilter] = React.useState({
    ...defaultFilter,
    ...query,
  })

  const { data: items, error } = useSWR([CATALOGS, filter], fetcher)
  React.useEffect(() => {
    setFilter({ ...filter, ...query })
  }, [query])

  const routerPush = f => {
    router.push(`/search?${querystring.stringify(f)}`)
  }

  const handleSearchSubmit = q => {
    const f = { ...filter, q, page: 1 }
    setFilter(f)
    routerPush({ q })
  }
  const handleSearchClear = () => {
    setFilter({ ...filter, q: '' })
    routerPush()
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
          <SearchInput value={filter.q} onSubmit={handleSearchSubmit} onClear={handleSearchClear} />
          <br />

          {error && <div>failed to load</div>}
          {!items && <LinearProgress color="secondary" />}
          {!error && items && (
            <div>
              <ItemList items={items.data} />
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
