import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'

const useStyles = makeStyles(theme => ({}))

export default function Home() {
  const classes = useStyles()

  return (
    <div className="container">
      <Header />

      <main>
        <Container maxWidth="md">
          <Typography variant="h2">
            <strong>Dota 2 Giftables</strong>
          </Typography>
          <TextField placeholder="Search" variant="outlined" />
        </Container>
      </main>

      <Footer />
    </div>
  )
}
