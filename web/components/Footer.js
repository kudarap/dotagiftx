import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import MuiLink from '@material-ui/core/Link'
import Divider from '@material-ui/core/Divider'
import Container from '@/components/Container'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  root: {
    marginTop: theme.spacing(10),
    marginBottom: theme.spacing(2),
  },
  list: {
    [theme.breakpoints.down('sm')]: {
      padding: theme.spacing(0, 1, 0),
    },
    [theme.breakpoints.down('xs')]: {
      display: 'flex',
      justifyContent: 'space-evenly',
      flexWrap: 'wrap',
    },
    display: 'block',
    listStyle: 'none',
    padding: 0,
    margin: theme.spacing(1, 0, 0),
    '& li': {
      [theme.breakpoints.down('xs')]: {
        float: 'none',
        marginRight: 0,
      },
      float: 'left',
      marginRight: theme.spacing(2),
    },
  },
}))

export default function () {
  const classes = useStyles()

  return (
    <footer className={classes.root}>
      <Divider />
      <Container disableMinHeight>
        <ul className={classes.list}>
          <li>
            <Link href="/about">About</Link>
          </li>
          <li>
            <Link href="/faq">FAQ</Link>
          </li>
          <li>
            <Link href="/privacy">Privacy</Link>
          </li>
          <li>
            <MuiLink
              href="http://vercel.com"
              target="_blank"
              color="textSecondary"
              rel="noreferrer noopener">
              Powered by Vercel
            </MuiLink>
          </li>
          {/* <li> */}
          {/*  <MuiLink href="http://chiligarlic.com" target="_blank"> */}
          {/*    A chiliGarlic project */}
          {/*  </MuiLink> */}
          {/* </li> */}
        </ul>
      </Container>
    </footer>
  )
}
