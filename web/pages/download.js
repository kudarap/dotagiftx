import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import MuiLink from '@mui/material/Link'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'
import Footer from '@/components/Footer'
import Button from '@/components/Button'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
  question: {
    paddingTop: theme.spacing(2.5),
    '&:target': {
      borderBottom: `2px inset ${theme.palette.secondary.main}`,
      '& .MuiLink-root:hover': {
        textDecoration: 'none',
      },
    },
  },
}))

function slugify(s) {
  return String(s)
    .toLowerCase()
    .replace(/[^a-z0-9 -]/g, '')
    .replace(/\s+/g, '-')
}

export function Title({ children, ...other }) {
  const { classes } = useStyles()
  const id = slugify(children)
  return (
    <Typography
      className={classes.question}
      component="h2"
      id={id}
      gutterBottom
      style={{ fontWeight: 'bold' }}
      {...other}>
      <MuiLink href={`#${id}`} color="textPrimary">
        {children}
      </MuiLink>
    </Typography>
  )
}

export default function Faqs() {
  const { classes } = useStyles()

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: DotagiftX for Mobile</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography
            sx={{ mt: 8, mb: 1 }}
            variant="h3"
            component="h1"
            fontWeight="bold"
            color="pimary">
            DotagiftX for Mobile
          </Typography>

          <Typography color="textSecondary" variant="h6" sx={{ mb: 2 }}>
            Open source mobile app for DotagiftX, You can checkout the repository on {` `}
            <Link
              href="https://github.com/tentenponce/dotagiftx-mobile"
              target="_blank"
              rel="noreferrer noopener">
              https://github.com/tentenponce/dotagiftx-mobile
            </Link>
            {` `} and you can download the latest release on:
          </Typography>

          <Title>1. Google Playstore</Title>
          <Typography color="textSecondary" gutterBottom>
            https://play.google.com/store/apps/details?id=com.dotagiftx
          </Typography>
          <Link
            href="https://play.google.com/store/apps/details?id=com.dotagiftx"
            target="_blank"
            rel="noreferrer noopener">
            <Button size="large" variant="contained" color="primary">
              Android
            </Button>
          </Link>

          <Title>2. Github Releases</Title>
          <Typography color="textSecondary" gutterBottom>
            https://github.com/tentenponce/dotagiftx-mobile/releases
          </Typography>
          <Link
            href="https://github.com/tentenponce/dotagiftx-mobile/releases"
            target="_blank"
            rel="noreferrer noopener">
            <Button size="large" variant="contained" color="primary">
              Build Release
            </Button>
          </Link>

          <Typography color="textSecondary" sx={{ mt: 6, mb: 12 }}>
            DotagiftX mobile has limited functionality compared to website and currently in active
            development. You can send feedback to{' '}
            <Link
              href="https://github.com/tentenponce/dotagiftx-mobile/issues"
              target="_blank"
              rel="noreferrer noopener">
              https://github.com/tentenponce/dotagiftx-mobile/issues
            </Link>
            .
          </Typography>
        </Container>
      </main>

      <Footer />
    </>
  )
}
