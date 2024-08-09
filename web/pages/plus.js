import React, { useContext } from 'react'
import Head from 'next/head'
import { USER_SUBSCRIPTION_MAP_COLOR, USER_SUBSCRIPTION_MAP_LABEL } from '@/constants/user'
import AppContext from '@/components/AppContext'
import { myProfile } from '@/service/api'
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
import Link from '@/components/Link'
import { Alert } from '@mui/material'
import { dateCalendar } from '@/lib/format'

const FeatureList = styled('ul')(({ theme }) => ({
  listStyle: 'none',
  '& li:before': {
    content: `'✔'`,
    marginRight: 8,
  },
  paddingLeft: 0,
}))

const defaultProfile = {
  subscription: null,
  subscription_type: '',
  subscription_ends_at: null,
  // runtime props
  subscriptionLabel: '',
  subscriptionColor: '',
}

export default function Plus() {
  const { isLoggedIn } = useContext(AppContext)

  // load subscription data if logged in.
  const [profile, setProfile] = React.useState(defaultProfile)
  React.useEffect(() => {
    ;(async () => {
      if (!isLoggedIn) {
        return
      }

      const res = await myProfile.GET()
      setProfile({
        ...res,
        subscriptionLabel: USER_SUBSCRIPTION_MAP_LABEL[res.subscription],
        subscriptionColor: USER_SUBSCRIPTION_MAP_COLOR[res.subscription],
      })
    })()
  }, [])

  return (
    <div className="container">
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Plus</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <Box>
            <Typography
              sx={{ mt: 8 }}
              style={{
                background: 'linear-gradient( to right, #CB8F37 20%, #F0CF59 50%, #B5793D 80% )',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                letterSpacing: 3,
                // textTransform: 'uppercase',
              }}
              variant="h3"
              component="h1"
              fontWeight="bold"
              color="secondary">
              Dotagift Plus
            </Typography>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Help support DotagiftX and get exclusive feature access, dedicated support, and
              profile badge.
            </Typography>
            {profile.subscription ? (
              <Alert
                variant="filled"
                sx={{
                  mt: 1,
                  background: profile.subscriptionColor,
                  transition: `box-shadow .5s ease-in-out, border .2s`,
                  borderTop: `0 solid ${profile.subscriptionColor}`,
                  boxShadow: `0 0 10px ${profile.subscriptionColor}`,
                }}>
                <strong>
                  You have an active {profile.subscriptionLabel} subscription in{' '}
                  {profile.subscription_type}{' '}
                  {profile.subscription_ends_at
                    ? `ends on ${dateCalendar(profile.subscription_ends_at)}`
                    : 'automatically'}
                  .
                </strong>
              </Alert>
            ) : null}
          </Box>

          <Grid container spacing={2} sx={{ mt: 0 }}>
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
                  <li>Supporter Badge</li>
                  <li>
                    <Link href="#exclusive-features">Refresher Shard</Link>
                  </li>
                </Box>
                <Button
                  variant="outlined"
                  fullWidth
                  sx={{ bgcolor: 'rgb(78, 93, 128)' }}
                  component={Link}
                  href="/transmute/subscription?id=supporter">
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
                  <li>Trader Badge</li>
                  <li>
                    <Link href="#exclusive-features">Refresher Orb</Link>
                  </li>
                </Box>
                <Button
                  variant="outlined"
                  fullWidth
                  sx={{ bgcolor: 'rgb(100, 159, 192)' }}
                  component={Link}
                  href="/transmute/subscription?id=trader">
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
                  <li>Partner Badge</li>
                  <li>
                    <Link href="#exclusive-features">Refresher Orb</Link>
                  </li>
                  <li>
                    <Link href="#exclusive-features">Shopkeeper's Contract</Link>
                  </li>
                  <li>
                    <Link href="#exclusive-features">Dedicated Pos-5</Link>
                  </li>
                </Box>
                <Button
                  variant="outlined"
                  fullWidth
                  sx={{ bgcolor: 'rgb(197, 144, 35)' }}
                  component={Link}
                  href="/transmute/subscription?id=partner">
                  <Typography variant="h6" sx={{ mr: 0.2 }}>
                    $20
                  </Typography>
                  /mo
                </Button>
              </Box>
            </Grid>
          </Grid>
          <br />
          <Typography variant="body2" paragraph textAlign="center" color="text.secondary">
            Subscriptions automatically renew and you can cancel your subscription on Paypal
            dashboard.
          </Typography>

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }} id="exclusive-features">
              Exclusive Features
            </Typography>
            <Grid container spacing={1.5}>
              {/* <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 4 }}>
                  <img src="/assets/badge-bp.png" height={48} />
                  <Typography>Partner Badge</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Simple cosmetic enhancement on your profile
                  </Typography>
                </Box>
              </Grid> */}

              <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 4 }}>
                  <img src="/assets/refresher-shard.png" height={48} />
                  <Typography>Refresher Shard</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Automatically refreshes expiring buy orders
                  </Typography>
                </Box>
              </Grid>

              <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 4 }}>
                  <img src="/assets/refresher-orb.png" height={48} />
                  <Typography>Refresher Orb</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Automatically refreshes expiring buy orders and listings
                  </Typography>
                </Box>
              </Grid>

              <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 4 }}>
                  <img src="/assets/recipe.png" height={48} />
                  <Typography>Shopkeeper's Contract</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Grants the ability to resell items outside your inventory
                  </Typography>
                </Box>
              </Grid>

              <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 4 }}>
                  <img src="/assets/courier.png" height={48} />
                  <Typography>Dedicated Pos-5</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Exclusive support channel on Discord and Steam
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          </Box>

          <Box
            sx={{
              display: 'none',
              mt: 8,
              p: 4,
              textAlign: 'center',
              background: 'url(/assets/plus-banner.png) no-repeat top center',
            }}>
            <Typography variant="h6">Partnership Goals</Typography>
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
          <TimelineOppositeContent color="text.secondary">5 subscribers</TimelineOppositeContent>
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
