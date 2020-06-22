import React from 'react'
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

export default function Search() {
  const classes = useStyles()

  const router = useRouter()
  const { q: keyword } = router.query

  const { data, error } = useSWR(`/items?q=${keyword}&sort=name`, fetcher)

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <SearchInput value={keyword} />

          <br />

          {error && <div>failed to load</div>}
          {!data && <div>loading...</div>}
          {data && <ItemList result={data} />}
        </Container>
      </main>

      <Footer />
    </>
  )
}
