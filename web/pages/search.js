import React from 'react'
import useSWR from 'swr'
import querystring from 'querystring'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import ItemList from '@/components/ItemList'
import SearchInput from '@/components/SearchInput'
import TablePagination from '@/components/TablePagination'
import { ITEMS, fetcher } from '@/service/api'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(2.5),
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

  const [filter, setFilter] = React.useState({ ...defaultFilter, page: query.page })
  const { data: items, error } = useSWR([ITEMS, filter], fetcher)

  React.useEffect(() => {
    const { page = 1 } = query
    setFilter({ ...filter, page })
  }, [query])

  const handleSearchChange = q => {
    setFilter({ ...filter, q })
  }

  const handleChangePage = (e, page) => {
    const f = { ...filter, page }
    setFilter(f)
    router.push(`/search?${querystring.stringify(f)}`)
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <SearchInput value={filter.q} onChange={handleSearchChange} />

          <br />

          {error && <div>failed to load</div>}
          {!items && <div>loading...</div>}
          {!error && items && (
            <div>
              <ItemList items={items.data} onChangePage={handleChangePage} />
              <TablePagination
                colSpan={3}
                style={{ textAlign: 'right' }}
                page={filter.page}
                count={items.total_count}
                onChangePage={handleChangePage}
              />
            </div>
          )}
        </Container>
      </main>

      <Footer />
    </>
  )
}
