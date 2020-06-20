import React from 'react'
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
  listControl: {},
  paginator: {
    float: 'right',
  },
}))

export default function Faq() {
  const classes = useStyles()

  const router = useRouter()
  const { q: keyword } = router.query

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <SearchInput value={keyword} />

          <div className={classes.listControl}>
            <Typography component="span">
              Results for &quot;<strong>{keyword}</strong>&quot;
            </Typography>
            <Typography component="span" className={classes.paginator}>
              Prev 1 2 3 4 5 Next
            </Typography>
          </div>

          <ItemList />
        </Container>
      </main>

      <Footer />
    </>
  )
}
