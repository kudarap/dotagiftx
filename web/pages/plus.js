import React from 'react'
import Head from 'next/head'
import { styled } from '@mui/material/styles'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import Timeline from '@mui/lab/Timeline'
import TimelineItem from '@mui/lab/TimelineItem'
import TimelineSeparator from '@mui/lab/TimelineSeparator'
import TimelineConnector from '@mui/lab/TimelineConnector'
import TimelineContent from '@mui/lab/TimelineContent'
import TimelineDot from '@mui/lab/TimelineDot'
import TimelineOppositeContent from '@mui/lab/TimelineOppositeContent'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { APP_NAME } from '@/constants/strings'
import tangoImage from '../public/assets/tango.png'

const FeatureList = styled('ul')(({ theme }) => ({
  listStyle: 'none',
  '& li:before': {
    content: `'✅ '`,
  },
  paddingLeft: theme.spacing(1),
}))

export default function Plus() {
  return (
    <div className="container">
      <Head>
        <title>{APP_NAME} :: Plus</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <Box>
            <Typography
              variant="h3"
              component="h1"
              sx={{ mt: 8 }}
              fontWeight="bold"
              color="secondary">
              Dotagift Plus
            </Typography>
            <Typography variant="h6">
              Help support DotagiftX and get exclusive feature access, dedicated support, and
              profile cosmetic.
            </Typography>
          </Box>

          <Grid container spacing={4} sx={{ mt: 0 }}>
            <Grid item xs={12} sm={6} md={4}>
              <Box
                sx={{
                  py: 2,
                  px: 3,
                  borderTop: '2px solid #596b95',
                  backgroundImage: 'linear-gradient(#4654755c, #465475)',
                }}>
                <Typography variant="h6">Supporter</Typography>
                <Box component={FeatureList} sx={{ height: 97 }}>
                  <li>Profile Badge</li>
                  <li>Refresher Shard</li>
                </Box>
                <Button variant="outlined" fullWidth sx={{ bgcolor: 'rgb(78, 93, 128)' }}>
                  <Typography variant="h6" sx={{ mr: 0.2 }}>
                    $1
                  </Typography>
                  /mo
                </Button>
              </Box>
            </Grid>

            <Grid item xs={12} sm={6} md={4}>
              <Box
                sx={{
                  py: 2,
                  px: 3,
                  borderTop: '2px solid #629cbd',
                  backgroundImage: 'linear-gradient(#578ba863, #578ba8)',
                }}>
                <Typography variant="h6">Trader</Typography>
                <Box component={FeatureList} sx={{ height: 97 }}>
                  <li>Profile Badge</li>
                  <li>Refresher Orb</li>
                </Box>
                <Button variant="outlined" fullWidth sx={{ bgcolor: 'rgb(100, 159, 192)' }}>
                  <Typography variant="h6" sx={{ mr: 0.2 }}>
                    $3
                  </Typography>
                  /mo
                </Button>
              </Box>
            </Grid>

            <Grid item xs={12} sm={12} md={4}>
              <Box
                sx={{
                  py: 2,
                  px: 3,
                  borderTop: '2px solid #ae7f1e',
                  backgroundImage: 'linear-gradient(#a6791d63, #a6791d)',
                  maxWidth: 500,
                  margin: 'auto',
                }}>
                <Typography variant="h6">Partner</Typography>
                <Box component={FeatureList}>
                  <li>Profile Badge</li>
                  <li>Refresher Orb</li>
                  <li>Shopkeeper's Contract</li>
                  <li>Dedicated Pos-5</li>
                </Box>
                <Button variant="outlined" fullWidth sx={{ bgcolor: 'rgb(197, 144, 35)' }}>
                  <Typography variant="h6" sx={{ mr: 0.2 }}>
                    $20
                  </Typography>
                  /mo
                </Button>
              </Box>
            </Grid>
          </Grid>
          <br />
          <Typography variant="caption">Subscriptions automatically renew</Typography>

          <Box sx={{ mt: 5 }}>
            <Typography variant="h6">Exclusive Features</Typography>
            <ul>
              <li>Partner Badge</li>
              <li>Refresher Shard - Automatically refreshes expiring buy orders.</li>
              <li>Refresher Orb - Automatically refreshes expiring buy orders and listings.</li>
              <li>
                Shopkeeper's Contract - Grants the ability to resell items outside your inventory.
              </li>
              <li>Dedicated Pos-5 - Exclusive support channel on Discord and Steam.</li>
            </ul>
          </Box>

          <Box
            sx={{
              mt: 10,
              p: 4,
              textAlign: 'center',
              background: 'url(/assets/plus-banner.png) no-repeat top center',
            }}>
            <Typography variant="h6">Unlockable Features</Typography>
            <FeatureUnlockables />
            <Typography variant="caption">
              Locked features will start development when goal is reached.
            </Typography>
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}

function FeatureUnlockables() {
  return (
    <React.Fragment>
      <Timeline>
        <TimelineItem>
          <TimelineOppositeContent color="text.secondary">3 subscribers</TimelineOppositeContent>
          <TimelineSeparator>
            <TimelineDot />
            <TimelineConnector />
          </TimelineSeparator>
          <TimelineContent>
            <Typography>Gem of Truesight</Typography>
            <Typography variant="body2" color="text.secondary">
              Grants vision to all buy orders
            </Typography>
          </TimelineContent>
        </TimelineItem>
        <TimelineItem>
          <TimelineOppositeContent color="text.secondary">10 subscribers</TimelineOppositeContent>
          <TimelineSeparator>
            <TimelineDot />
            <TimelineConnector />
          </TimelineSeparator>
          <TimelineContent>
            <Typography>Seer Stone</Typography>
            <Typography variant="body2" color="text.secondary">
              Provides an analytics and monitoring dashboard
            </Typography>
          </TimelineContent>
        </TimelineItem>
        <TimelineItem>
          <TimelineOppositeContent color="text.secondary">??</TimelineOppositeContent>
          <TimelineSeparator>
            <TimelineDot />
            <TimelineConnector />
          </TimelineSeparator>
          <TimelineContent>
            <Typography>Fusion Rune</Typography>
            <Typography variant="body2" color="text.secondary">
              Ability to create your own cache
            </Typography>
          </TimelineContent>
        </TimelineItem>
        <TimelineItem>
          <TimelineOppositeContent color="text.secondary">???</TimelineOppositeContent>
          <TimelineSeparator>
            <TimelineDot />
          </TimelineSeparator>
          <TimelineContent>???</TimelineContent>
        </TimelineItem>
      </Timeline>
    </React.Fragment>
  )
}
