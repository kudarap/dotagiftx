import React from 'react'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Container from '@/components/Container'
import Header from '@/components/Header'
import Footer from '@/components/Footer'

const useStyles = makeStyles()(theme => ({
  main: {
    marginTop: theme.spacing(6),
  },
}))

// eslint-disable-next-line react/prop-types
function Error({ statusCode }) {
  const { classes } = useStyles()

  return (
    <>
      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom align="center">
            Internal Error
          </Typography>

          <Typography color="textSecondary" align="center">
            {statusCode
              ? `An error ${statusCode} occurred on server`
              : 'An error occurred on client'}
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}

Error.getInitialProps = ({ res, err }) => {
  // eslint-disable-next-line no-nested-ternary
  const statusCode = res ? res.statusCode : err ? err.statusCode : 404
  return { statusCode }
}

export default Error
