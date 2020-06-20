import React from 'react'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(10),
  },
}))

export default function Faq() {
  const classes = useStyles()

  const router = useRouter()
  const { slug } = router.query

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography>
            Item page for <strong>{slug}</strong>!
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
