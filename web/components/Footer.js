import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'
import Container from '@/components/Container'
import Link from '@/components/Link'
import { APP_NAME } from '@/constants/strings'
import { APP_FOOTER_HEIGHT } from '@/constants/app'

// Stole from SteamDB dota 2 page footer.
const dotaHeroImage =
  'https://cdn.cloudflare.steamstatic.com/steam/apps/570/library_hero.jpg?t=1605830961'

const useStyles = makeStyles(theme => ({
  root: {
    [theme.breakpoints.down('xs')]: {
      // Keeps the footer on the bottom of the screen on small screens.
      height: APP_FOOTER_HEIGHT + 10,
    },
    [theme.breakpoints.down('sm')]: {
      paddingBottom: theme.spacing(0),
    },
    marginTop: theme.spacing(5),
    paddingBottom: theme.spacing(5),
    height: APP_FOOTER_HEIGHT,
    background: `linear-gradient(to bottom, rgba(38, 50, 56, 0.7) ${
      APP_FOOTER_HEIGHT + 9
    }px, transparent), url(${dotaHeroImage}) center -140px`,
  },
  list: {
    [theme.breakpoints.down('xs')]: {
      display: 'flex',
      justifyContent: 'space-evenly',
      flexWrap: 'wrap',
      margin: 0,
    },
    display: 'block',
    listStyle: 'none',
    padding: 0,
    '& li': {
      [theme.breakpoints.down('xs')]: {
        float: 'none',
        marginRight: 0,
        paddingTop: theme.spacing(1),
        paddingBottom: theme.spacing(1),
      },
      float: 'left',
      marginRight: theme.spacing(2),
    },
  },
  vavleCopyright: {
    [theme.breakpoints.down('xs')]: {
      marginTop: theme.spacing(0),
      marginBottom: theme.spacing(3),
      textAlign: 'center',
    },
    display: 'block',
    marginTop: theme.spacing(4),
  },
}))

export default function Footer() {
  const classes = useStyles()

  return (
    <footer className={classes.root}>
      {/* <Container disableMinHeight> */}
      {/*  <Typography variant="caption" color="textSecondary"> */}
      {/*    Game content and materials are trademarks and copyrights of their respective publisher and */}
      {/*    its licensors. All rights reserved. */}
      {/*  </Typography> */}
      {/* </Container> */}
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
            <Link href="/donate" color="textSecondary">
              Donate
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
                {APP_NAME}
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
        <br />
        <Typography
          className={classes.vavleCopyright}
          variant="caption"
          color="textSecondary"
          component="p">
          {APP_NAME} is a community website and not affiliated with Valve or Steam.
          <br />
          Game content and materials are trademarks and copyrights of their respective publisher and
          its licensors. All rights reserved.
        </Typography>
      </Container>
    </footer>
  )
}
