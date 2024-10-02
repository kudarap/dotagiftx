import React from 'react'
import { makeStyles } from 'tss-react/mui'
import Divider from '@mui/material/Divider'
import Typography from '@mui/material/Typography'
import Container from '@/components/Container'
import Link from '@/components/Link'
import { APP_NAME } from '@/constants/strings'
import { APP_FOOTER_HEIGHT } from '@/constants/app'

// Stole from SteamDB dota 2 page footer.
const t = 1724395576617
const dotaHeroImage = `https://cdn.cloudflare.steamstatic.com/steam/apps/570/library_hero.jpg?t=${t}`
// const dotaHeroImage =
;('https://clan.cloudflare.steamstatic.com/images/3703047/ba80108f618e691d184e7eb5579e56c33b9a811b.jpg')

// const heroImage = '/assets/bg_hero.png'

// const mobileHeightCompensator = 31
const mobileHeightCompensator = 100

const useStyles = makeStyles()(theme => ({
  root: {
    [theme.breakpoints.down('sm')]: {
      // Keeps the footer on the bottom of the screen on small screens.
      height: APP_FOOTER_HEIGHT + mobileHeightCompensator,
      // backgroundPositionX: '0, -70%',
    },
    [theme.breakpoints.down('md')]: {
      paddingBottom: theme.spacing(0),
    },
    marginTop: theme.spacing(5),
    height: APP_FOOTER_HEIGHT,
    background: `linear-gradient(0deg, rgba(38, 50, 56, 0.36) 0%, rgb(38, 50, 56) 100%), url(${dotaHeroImage}) center -50px`,
  },
  list: {
    [theme.breakpoints.down('sm')]: {
      display: 'flex',
      justifyContent: 'space-evenly',
      flexWrap: 'wrap',
      margin: 0,
    },
    display: 'block',
    listStyle: 'none',
    padding: 0,
    '& li': {
      [theme.breakpoints.down('sm')]: {
        float: 'none',
        marginRight: theme.spacing(1),
        marginLeft: theme.spacing(1),
        paddingTop: theme.spacing(1),
        paddingBottom: theme.spacing(1),
      },
      float: 'left',
      marginRight: theme.spacing(2),
    },
  },
  vavleCopyright: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(0),
      marginBottom: theme.spacing(3),
      textAlign: 'center',
    },
    display: 'block',
    marginTop: theme.spacing(4),
  },
  highlight: {
    // background: '-webkit-linear-gradient(#EBCF87 10%, #C79123 90%)',
    background: '-webkit-linear-gradient(#EBCF87 10%, #A9EFAA, #7FBC8B)',
    backgroundClip: 'border-box',
    backgroundClip: 'text',
    WebkitBackgroundClip: 'text',
    WebkitTextFillColor: 'transparent',
    filter: 'drop-shadow(0px 0px 5px #e1261c)',
  },
}))

function SteamAwardIcon(props) {
  return (
    <svg viewBox="0 0 13 21" xmlns="http://www.w3.org/2000/svg" aria-hidden="true" {...props}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M5.8.2l-.6.42c-.24.15-.5.23-.81.2L3.93.73c-.5-.04-.97.19-1.2.65l-.3.65c-.16.27-.35.47-.62.58l-.42.2c-.46.19-.73.65-.7 1.15l.05.73c.03.3-.04.54-.24.77l-.26.39c-.31.38-.31.92-.04 1.34l.42.62c.15.23.2.5.2.8l-.08.47c-.04.5.19.96.65 1.2l.65.3c.27.15.46.35.58.62l.15.42c.23.46.66.73 1.2.7l.72-.05c.27-.04.54.04.77.23l.39.27c.38.31.92.31 1.34.04l.62-.42c.23-.16.5-.2.77-.2l.46.08c.5.04 1-.19 1.19-.65l.34-.66c.12-.27.31-.46.58-.57l.42-.16c.46-.23.77-.65.73-1.19l-.04-.73a.97.97 0 01.24-.77l.26-.39c.31-.38.31-.92.04-1.34l-.42-.62a1.1 1.1 0 01-.2-.77l.08-.46c.04-.5-.19-1-.65-1.2l-.65-.34c-.27-.12-.46-.3-.58-.58l-.2-.42C10 .97 9.55.66 9.05.7L8.3.74A.93.93 0 017.54.5L7.15.24C6.77-.07 6.23-.07 5.81.2zm.7 2.5a3.82 3.82 0 10-.04 7.65A3.82 3.82 0 006.5 2.7z"
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M3.04 12.78v7.32l3.46-2.47 3.46 2.47v-7.32c-.3.12-.65.2-1 .16l-.42-.08a.65.65 0 00-.39.12l-.61.38c-.66.42-1.46.42-2.12-.04l-.34-.3c-.12-.08-.23-.08-.35-.08l-.73.04c-.34 0-.69-.04-.96-.2z"
      />
    </svg>
  )
}

export default function Footer() {
  const { classes } = useStyles()

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
            <Link href="/faqs" color="textSecondary">
              FAQs
            </Link>
          </li>
          <li>
            <Link href="/middlemen" color="textSecondary">
              Middlemen
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
              <Typography
                variant="body2"
                component="span"
                color="secondary"
                className={classes.highlight}>
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
          <li style={{ float: 'right' }}>
            <Link
              href="https://steamcommunity.com/sharedfiles/filedetails/?id=2313234224"
              style={{ color: '#ffc83d' }}>
              <SteamAwardIcon
                style={{ margin: '0 2px -8px 0' }}
                fill="#ffc83d"
                width="24"
                height="24"
              />
              <strong>Give a Steam Award</strong>
            </Link>
          </li>
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
