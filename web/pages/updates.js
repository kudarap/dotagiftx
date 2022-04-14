import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Divider from '@mui/material/Divider'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

export default function Updates() {
  const { classes } = useStyles()

  return (
    <div className="container">
      <Head>
        <title>{APP_NAME} :: Updates</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" gutterBottom>
            April 10, 2022
          </Typography>
          {/* <Typography variant="body" gutterBottom>
            v0.18.0: Dotagift+
          </Typography> */}
          <Typography color="textSecondary">
            <ul>
              <li>added subscription page</li>
              <li>added treasures page</li>
              <li>added update page</li>
              <li>added search dialog with top queries</li>
              <li>reworked navigation header</li>
            </ul>
          </Typography>
          <br />
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            Feb 27, 2022
          </Typography>
          {/* <Typography variant="body" gutterBottom>
            v0.17.1: Enable SSR
          </Typography> */}
          <Typography color="textSecondary">
            <ul>
              <li>added inventory checker cli tool</li>
              <li>added lazy loading graph and stats on item page</li>
              <li>implemented emotion cache</li>
              <li>migrated styling from jss to tss</li>
              <li>updated ssr style cache using emotion and tss</li>
              <li>updated logo color with blue and yellow</li>
            </ul>
          </Typography>
          <br />
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            Feb 27, 2022
          </Typography>
          {/* <Typography variant="body" gutterBottom>
            v0.17.0: Improved logo and material-ui v5 migration
          </Typography> */}
          <Typography color="textSecondary">
            <ul>
              <li>added hammer service endpoint</li>
              <li>added default read and write server timeout</li>
              <li>added rethinkdb timeout, cap and max connection</li>
              <li>added production remote dump and sync script</li>
              <li>added new static pages and emoved user profiles on sitemap</li>
              <li>added dynamic sitemap from generated api sitemap</li>
              <li>updated and migrated mui from v4 to v5</li>
              <li>updated logo and branding</li>
              <li>updated footer bg for primal beast update</li>
              <li>excluded profiles pages on crawlers</li>
              <li>fixes comma on buy orders total count</li>
            </ul>
          </Typography>
          <br />
          <Divider />
          <br />
        </Container>
      </main>

      <Footer />
    </div>
  )
}
