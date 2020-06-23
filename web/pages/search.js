import React from 'react'
import querystring from 'querystring'
import useSWR from 'swr'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import ItemList from '@/components/ItemList'
import SearchInput from '@/components/SearchInput'
import { fetcher } from '../service/api'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(2.5),
  },
}))

const defaultFilter = {
  sort: 'name',
}

function parseParams(url, filter) {
  console.log(`parseParams ${url}?${querystring.stringify(filter)}`)
  // return url
  return `${url}?${querystring.stringify(filter)}`
}

export default function Search() {
  const classes = useStyles()

  const router = useRouter()
  const { q } = router.query

  const [filter, setFilter] = React.useState({ ...defaultFilter, q })

  // const { data, error } = useSWR(['/items', filter], fetcher)
  const { data, error } = useSWR(parseParams('/items', { ...filter, q }), fetcher)

  const handleChangePage = (e, page) => {
    console.log(page)
  }

  const handleSearchChange = value => {
    setFilter({ ...filter, q: value })
  }

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <SearchInput value={filter.q} onChange={handleSearchChange} />

          <br />

          {error && <div>failed to load</div>}
          {!data && <div>loading...</div>}
          {!error && data && <ItemList result={data} onChangePage={handleChangePage} />}
        </Container>
      </main>

      <Footer />
    </>
  )
}
