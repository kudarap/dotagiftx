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
    marginBottom: theme.spacing(4),
  },
  banner: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(0),
    },
    margin: theme.spacing(20, 0, 4, 0),
  },
  bannerText: {
    [theme.breakpoints.down('sm')]: {
      fontSize: 35,
    },
    fontWeight: 'bold',
  },
}))

function Banner() {
  const classes = useStyles()

  return (
    <div className={classes.banner}>
      <Typography className={classes.bannerText} variant="h3" align="center">
        Search Dota 2 <span style={{ display: 'inline-block' }}>giftable items</span>
      </Typography>
    </div>
  )
}

export default function Home() {
  const classes = useStyles()

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Banner />

          <TextField
            className={classes.searchBar}
            fullWidth
            placeholder="Search Item, Hero, Treasure..."
            helperText="search on 332 for posted items"
            variant="outlined"
            color="secondary"
          />

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
