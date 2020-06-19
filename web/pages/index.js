import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(10),
  },
  searchBar: {
    maxWidth: 640,
    margin: '0 auto',
    display: 'block',
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
            variant="outlined"
            color="secondary"
          />
        </Container>
      </main>

      <Footer />
    </>
  )
}
