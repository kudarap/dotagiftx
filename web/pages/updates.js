import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Divider from '@mui/material/Divider'
import Link from '@/components/Link'
import LaunchIcon from '@mui/icons-material/Launch'

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
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Updates</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" gutterBottom>
            Jul 21, 2024
          </Typography>
          <Link
            target="_blank"
            rel="noreferrer noopener"
            href="https://github.com/kudarap/dotagiftx/releases/tag/v0.21.0">
            {' '}
            v0.21.0: Auto Subscription Crackdown and Crownfall Treasure{' '}
            <LaunchIcon fontSize="relative" />
          </Link>
          <br />
          <br />
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            May 6, 2024
          </Typography>
          <Link
            target="_blank"
            rel="noreferrer noopener"
            href="https://github.com/kudarap/dotagiftx/releases/tag/v0.20.0">
            {' '}
            v0.20.0: Span tracing, Task Queue, and Optimizations <LaunchIcon fontSize="relative" />
          </Link>
          <br />
          <br />
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            Sep 27, 2023
          </Typography>
          <Link
            target="_blank"
            rel="noreferrer noopener"
            href="https://github.com/kudarap/dotagiftx/releases/tag/v0.19.0">
            {' '}
            v0.19.0: Task verification + bunch of fixes and optimizations{' '}
            <LaunchIcon fontSize="relative" />
          </Link>
          <br />
          <br />
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            October 8, 2022
          </Typography>
          <Typography color="textSecondary">
            <ul>
              <li>added immortal treasure 1 2022</li>
              <li>changes vanity to redirect automaticaly to profile id</li>
              <li>changes listing quantity limit to 1, user has refresher orb limit to 5</li>
              <li>updated donate page copy that donator badge is not available anymore</li>
              <li>updated theme for TI 2022 event</li>
            </ul>
          </Typography>
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            September 3, 2022
          </Typography>
          <Typography color="textSecondary">
            <ul>
              <li>added expired market sweeper</li>
              <li>added optimization on catalog indexing and invalidation</li>
              <li>added 3 points for resell deliveries on user score</li>
              <li>added refresher shard boon on trader subscription</li>
              <li>added giveaway link</li>
              <li>added middleman badge</li>
              <li>fixes banned profile display</li>
              <li>removed expiring post notice</li>
            </ul>
          </Typography>
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            April 15, 2022
          </Typography>
          {/* <Typography variant="body" gutterBottom>
            v0.18.0: Dotagift+
          </Typography> */}
          <Typography color="textSecondary">
            <ul>
              <li>added subscription page</li>
              <li>added treasures page</li>
              <li>added update page</li>
              <li>added expiring post page</li>
              <li>added search dialog with top queries</li>
              <li>reworked navigation header</li>
            </ul>
          </Typography>
          <Typography>
            Read more about{' '}
            <Link href="/expiring-posts" color="secondary">
              Expiring Posts
            </Link>
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
