import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import { APP_NAME } from '@/constants/strings'
import Container from '@/components/Container'
import Header from '@/components/Header'
import Footer from '@/components/Footer'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function Custom404({ children }) {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: 404 - Page Not Found</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          {children || (
            <Typography variant="h5" component="h1" gutterBottom align="center">
              404 - Page Not Found
            </Typography>
          )}
        </Container>
      </main>

      <Footer />
    </>
  )
}
