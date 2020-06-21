import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import ItemListPopular from '@/components/ItemListPopular'
import ItemListRecent from '@/components/ItemListRecent'
import SearchInput from '@/components/SearchInput'

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
    color: theme.palette.app.white,
  },
}))

function Banner() {
  const classes = useStyles()

  return (
    <div className={classes.banner}>
      <Typography className={classes.bannerText} variant="h3" align="center">
        Search for Dota 2 <span style={{ display: 'inline-block' }}>giftable items</span>
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

          <SearchInput helperText="Search on 92 for sale items" />

          <br />
          <Typography>Popular Items</Typography>
          <ItemListPopular />
          <br />
          <Typography>Recently Posted</Typography>

          <ItemListRecent />
        </Container>
      </main>

      <Footer />
    </>
  )
}
