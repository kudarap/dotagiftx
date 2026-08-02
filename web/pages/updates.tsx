import React from 'react'
import Head from 'next/head'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Divider from '@mui/material/Divider'
import LaunchIcon from '@mui/icons-material/Launch'
import { APP_NAME } from '@/constants/strings'
import Footer from '@/components/Footer'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
  },
}))

interface Release {
  releaseDate: string
  tag: string
  label: string
  contents: string[]
}

const releaseItem = (releaseDate: string, tag: string, label: string, contents: string[]): Release => ({
  releaseDate,
  tag,
  label,
  contents,
})

const releases: Release[] = [
  releaseItem('Jun 25, 2026', 'v0.24.0', 'Dark Carnival', [
    'added dark carnival event theme',
    'added umami web analytics in deprecation of google analytics',
    'added new tab on my listing page to fix broken steam id',
    'added steam id on profile page and card',
    'added phantasm support for multiple cloud providers',
    'updated middleman page dynamically load users',
    'updated pipeline for new infrastructure',
    'improved phantasm verification retries',
    'migrated caching from redis to valkey',
    'fixed phantasm always using the fallback inventory provider',
    'fixed dev env not showing images',
    'fixed search dialog removing spaces',
    'fixed vanity resolver prevent user to search by stean vanity id',
    'fixed nextjs build warnings',
    'fixed docker image vulnerability',
    'fixed buyer orders redaction',
    'fixed config requiring .env file to exists',
  ]),
  releaseItem(
    'March 10, 2026',
    'v0.23.2',
    'Clickhouse integration, new pages and code improvements',
    [
      'migrated to eslint v6 to v8',
      'migrated from yarn v4.x to bun 1.x',
      'fixed some lints and some set to warnings',
      'fixed avatar broken image',
      'added phantasm cache cleaner',
      'updated vercel bun runtime from node to bun',
    ]
  ),
  releaseItem('Feb 22, 2026', 'v0.23.1', 'Housekeeping update: go v1.26 and yarn v4', [
    'upgraded to go1.26',
    'upgraded to yarn v4',
    'upgraded github actions',
    'updated go modules',
    'updated ui packages',
    'updated dockerfile for inline caching mount',
    'updated readme for local setup and requirements',
    'fixed lints on clickhouse',
  ]),
  releaseItem(
    'Feb 22, 2026',
    'v0.23.0',
    'Clickhouse integration, new pages and code improvements',
    [
      'added clickhouse integration',
      'added experimental market status endpoint with manual indexing',
      'added custom verified and elapsed time on item verification card',
      'added vercel middle ware as request middleware',
      'added request id on response header',
      'added TI14 theme and cosmic treasure',
      'added heroes page',
      'added mobile download page',
      'added treasures and heroes endpoints',
      'added hero and treasure endpoint',
      'added added winter 2025 collectors cache page',
      'added major data loss incident page',
      'added new hero largo',
      'fixed phantasm retry after ttl',
      'fixed immortal treasure 2020 page',
      'upgraded to go1.25',
      'reworked auth salt',
      'reworked caching',
      'updated middleman page',
      'optimized market creation validation query',
      'upgraded to new stats market endpoints',
      'flattened packages',
      'removed deprecated code',
      'removed stop killing games and update ui packages',
      'removed paypal module',
      'renamed custom errors to xerrors',
      'removed hash package',
      'removed xerror package',
    ]
  ),
]

export default function Updates() {
  const { classes } = useStyles()

  return (
    <div className="container">
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Updates</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Typography variant="h5" component="h1" gutterBottom style={{ marginBottom: 26 }}>
            Updates
          </Typography>

          {releases.map(release => (
            <div key={release.tag}>
              <Typography variant="h5" gutterBottom>
                {release.releaseDate}
              </Typography>
              <Link
                target="_blank"
                rel="noreferrer noopener"
                href={`https://github.com/kudarap/dotagiftx/releases/tag/${release.tag}`}>
                <Typography color="secondary" component="span">
                  {release.tag}
                </Typography>{' '}
                {release.label}
                <LaunchIcon fontSize="small" />
              </Link>
              <ul>
                {release.contents.map(line => (
                  <Typography key={line} color="textSecondary" component="li">
                    {line}
                  </Typography>
                ))}
              </ul>
              <Divider />
              <br />
            </div>
          ))}

          <Typography variant="h5" gutterBottom>
            June 25, 2025
          </Typography>
          <Link
            target="_blank"
            rel="noreferrer noopener"
            href="https://github.com/kudarap/dotagiftx/releases/tag/v0.22.0">
            {' '}
            v0.22.0: 5th year anniversary - Spring forward and Phantasm Crawl{' '}
            <LaunchIcon fontSize="small" />
          </Link>
          <br />
          <br />
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            Jul 21, 2024
          </Typography>
          <Link
            target="_blank"
            rel="noreferrer noopener"
            href="https://github.com/kudarap/dotagiftx/releases/tag/v0.21.0">
            {' '}
            v0.21.0: Auto Subscription Crackdown and Crownfall Treasure{' '}
            <LaunchIcon fontSize="small" />
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
            v0.20.0: Span tracing, Task Queue, and Optimizations <LaunchIcon fontSize="small" />
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
            <LaunchIcon fontSize="small" />
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
              <Typography color="textSecondary" component="li">
                added immortal treasure 1 2022
              </Typography>
              <Typography color="textSecondary" component="li">
                changes vanity to redirect automaticaly to profile id
              </Typography>
              <Typography color="textSecondary" component="li">
                changes listing quantity limit to 1, user has refresher orb limit to 5
              </Typography>
              <Typography color="textSecondary" component="li">
                updated donate page copy that donator badge is not available anymore
              </Typography>
              <Typography color="textSecondary" component="li">
                updated theme for TI 2022 event
              </Typography>
            </ul>
          </Typography>
          <Divider />
          <br />

          <Typography variant="h5" gutterBottom>
            September 3, 2022
          </Typography>
          <Typography color="textSecondary">
            <ul>
              <Typography color="textSecondary" component="li">
                added expired market sweeper
              </Typography>
              <Typography color="textSecondary" component="li">
                added optimization on catalog indexing and invalidation
              </Typography>
              <Typography color="textSecondary" component="li">
                added 3 points for resell deliveries on user score
              </Typography>
              <Typography color="textSecondary" component="li">
                added refresher shard boon on trader subscription
              </Typography>
              <Typography color="textSecondary" component="li">
                added giveaway link
              </Typography>
              <Typography color="textSecondary" component="li">
                added middleman badge
              </Typography>
              <Typography color="textSecondary" component="li">
                fixes banned profile display
              </Typography>
              <Typography color="textSecondary" component="li">
                removed expiring post notice
              </Typography>
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
              <Typography color="textSecondary" component="li">
                added subscription page
              </Typography>
              <Typography color="textSecondary" component="li">
                added treasures page
              </Typography>
              <Typography color="textSecondary" component="li">
                added update page
              </Typography>
              <Typography color="textSecondary" component="li">
                added expiring post page
              </Typography>
              <Typography color="textSecondary" component="li">
                added search dialog with top queries
              </Typography>
              <Typography color="textSecondary" component="li">
                reworked navigation header
              </Typography>
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
              <Typography color="textSecondary" component="li">
                added inventory checker cli tool
              </Typography>
              <Typography color="textSecondary" component="li">
                added lazy loading graph and stats on item page
              </Typography>
              <Typography color="textSecondary" component="li">
                implemented emotion cache
              </Typography>
              <Typography color="textSecondary" component="li">
                migrated styling from jss to tss
              </Typography>
              <Typography color="textSecondary" component="li">
                updated ssr style cache using emotion and tss
              </Typography>
              <Typography color="textSecondary" component="li">
                updated logo color with blue and yellow
              </Typography>
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
              <Typography color="textSecondary" component="li">
                added hammer service endpoint
              </Typography>
              <Typography color="textSecondary" component="li">
                added default read and write server timeout
              </Typography>
              <Typography color="textSecondary" component="li">
                added rethinkdb timeout, cap and max connection
              </Typography>
              <Typography color="textSecondary" component="li">
                added production remote dump and sync script
              </Typography>
              <Typography color="textSecondary" component="li">
                added new static pages and emoved user profiles on sitemap
              </Typography>
              <Typography color="textSecondary" component="li">
                added dynamic sitemap from generated api sitemap
              </Typography>
              <Typography color="textSecondary" component="li">
                updated and migrated mui from v4 to v5
              </Typography>
              <Typography color="textSecondary" component="li">
                updated logo and branding
              </Typography>
              <Typography color="textSecondary" component="li">
                updated footer bg for primal beast update
              </Typography>
              <Typography color="textSecondary" component="li">
                excluded profiles pages on crawlers
              </Typography>
              <Typography color="textSecondary" component="li">
                fixes comma on buy orders total count
              </Typography>
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
