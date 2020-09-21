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
            Who is behind this?
          </Typography>
          <br />

          <Avatar
            src="https://api.dotagiftx.com/images/adfb7fc8133861692abc5631d67b5f51dfd5753f.jpg"
            style={{ width: 100, height: 100, borderRadius: 0 }}
          />
          <Typography color="textSecondary">
            <strong>kudarap</strong> &mdash; programmer
            <br />
            Feel free to message me on{' '}
            <Link
              color="secondary"
              target="_blank"
              rel="noreferrer noopener"
              href="https://www.reddit.com/message/compose/?to=kudarap">
              Reddit
            </Link>{' '}
            if you have issues or suggestions.
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
