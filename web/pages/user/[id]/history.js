import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import { Link } from '@material-ui/core'

const useStyles = makeStyles(theme => ({
  main: {
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
          <Typography variant="h5" component="h1" gutterBottom>
            History
          </Typography>
          <br />
        </Container>
      </main>

      <Footer />
    </>
  )
}
