import React from 'react'
import fetch from 'unfetch'
import useSWR from 'swr'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import ItemList from '@/components/ItemList'
import SearchInput from '@/components/SearchInput'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(4),
  },
}))

const fetcher = url => fetch(url).then(r => r.json())

export default function Search() {
  const classes = useStyles()

  const router = useRouter()
  const { q: keyword } = router.query

  const { data, error } = useSWR(`http://192.168.1.79:8000/items?q=${keyword}&sort=name`, fetcher)

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <SearchInput value={keyword} />

          <Typography>
            Results for &quot;<strong>{keyword}</strong>&quot;
          </Typography>

          {error && <div>failed to load</div>}
          {!data && <div>loading...</div>}
          {data && <ItemList result={data} />}
        </Container>
      </main>

      <Footer />
    </>
  )
}
