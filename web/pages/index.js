import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import ItemList from '@/components/ItemList'
import ItemListRecent from '@/components/ItemListRecent'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(4),
  },
  searchBar: {
    margin: '0 auto',
  },
}))

export default function Home() {
  const classes = useStyles()

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <TextField
            className={classes.searchBar}
            fullWidth
            placeholder="Search Item, Hero, Treasure..."
            helperText="search on 332 for posted items"
            variant="outlined"
            color="secondary"
          />
          <br />
          <br />
          <Typography>Popular Items</Typography>
          <ItemList />
          <br />
          <Typography>Recently Posted</Typography>
          <ItemListRecent />
        </Container>
      </main>

      <Footer />
    </>
  )
}
