import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import MarketForm from '@/components/MarketForm'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(6),
  },
}))

export default function About() {
  const classes = useStyles()

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <MarketForm />
        </Container>
      </main>

      <Footer />
    </>
  )
}
