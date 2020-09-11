import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'
import Container from '@/components/Container'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  root: {
    [theme.breakpoints.down('sm')]: {
      paddingBottom: theme.spacing(0),
    },
    marginTop: theme.spacing(5),
    paddingBottom: theme.spacing(5),
  },
  list: {
    [theme.breakpoints.down('xs')]: {
      display: 'flex',
      justifyContent: 'space-evenly',
      flexWrap: 'wrap',
    },
    display: 'block',
    listStyle: 'none',
    padding: 0,
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

export default function Footer() {
  const classes = useStyles()

  return (
    <footer className={classes.root}>
      <Divider />
      <Container disableMinHeight>
        <ul className={classes.list}>
          {/* <li> */}
          {/*  <Link href="/about" color="textSecondary"> */}
          {/*    About */}
          {/*  </Link> */}
          {/* </li> */}
          <li>
            <Link href="/faq" color="textSecondary">
              FAQ
            </Link>
          </li>
          <li>
            <Link href="/privacy" color="textSecondary">
              Privacy
            </Link>
          </li>
          <li>
            {/* <MuiLink */}
            {/*  href="http://github.com/kudarap" */}
            {/*  target="_blank" */}
            {/*  color="textSecondary" */}
            {/*  rel="noreferrer noopener"> */}
            {/*  <Typography component="span" color="secondary"> */}
            {/*    DotagiftX */}
            {/*  </Typography>{' '} */}
            {/*  by kudarap */}
            {/* </MuiLink> */}
            <Link href="/about" color="textSecondary">
              <Typography component="span" color="secondary">
                DotagiftX
              </Typography>{' '}
              by kudarap
            </Link>
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
