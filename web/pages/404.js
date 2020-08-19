import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Container from '@/components/Container'
import Header from '@/components/Header'
import Footer from '@/components/Footer'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(6),
  },
}))

export default function Custom404() {
  const classes = useStyles()

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom align="center">
            404 - Page Not Found
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
