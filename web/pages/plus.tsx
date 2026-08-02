import React, { useContext } from 'react'
import Head from 'next/head'
import Image from 'next/image'
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
import { myProfile } from '@/service/api'
import AppContext from '@/components/AppContext'
import {
  USER_SUBSCRIPTION_MAP_COLOR,
  USER_SUBSCRIPTION_MAP_LABEL,
  USER_SUBSCRIPTION_PARTNER,
  USER_SUBSCRIPTION_SUPPORTER,
  USER_SUBSCRIPTION_TRADER,
} from '@/constants/user'
import { dateCalendar } from '@/lib/format'

const FeatureList = styled('ul')(() => ({
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

interface Profile {
  subscription: number | null
  subscription_type: string
  subscription_ends_at: string | null
  subscriptionLabel: string
  subscriptionColor: string
}

export default function Plus() {
  const { isLoggedIn } = useContext(AppContext)

  // load subscription data if logged in.
  const [profile, setProfile] = React.useState<Profile>(defaultProfile)
  React.useEffect(() => {
    ;(async () => {
      if (!isLoggedIn) {
        return
      }

      const res = (await myProfile.GET(true)) as {
        subscription?: number
        subscription_type?: string
        subscription_ends_at?: string | null
      }
      setProfile({
        ...(res as object),
        subscription: res.subscription ?? null,
        subscription_type: res.subscription_type || '',
        subscription_ends_at: res.subscription_ends_at || null,
        subscriptionLabel: USER_SUBSCRIPTION_MAP_LABEL[res.subscription || 0],
        subscriptionColor: USER_SUBSCRIPTION_MAP_COLOR[res.subscription || 0],
      })
    })()
  }, [isLoggedIn])

  const injectActiveStyle = (subscriptionId: number, style: { borderTop: string } & Record<string, unknown>) => {
    if (!profile.subscription || subscriptionId != profile.subscription) {
      return {
        ...style,
        borderColor: 'transparent',
      }
    }

    return {
      ...style,
      transition: `all .5s ease-in-out, border .2s`,
      MozTransition: `all .5s ease-in-out, border .2s`,
      boxShadow: `0 0 40px ${profile.subscriptionColor}`,
      border: style.borderTop,
      borderStyle: 'solid',
    }
  }

  return (
    <div className="container">
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Plus</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main>
        <Container>
          <Box>
            <Typography
              sx={{ mt: 8, mb: 1, letterSpacing: 3 }}
              style={{
                background: 'linear-gradient( to right, #CB8F37 20%, #F0CF59 50%, #B5793D 80% )',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                // textTransform: 'uppercase',
              }}
              variant="h3"
              component="h1"
              fontWeight="bold"
              color="secondary">
              Dotagift Plus
            </Typography>
            <Typography variant="h6" color="textSecondary" sx={{ mb: 2 }}>
              Help support the project and get exclusive feature access, <s>dedicated support</s>,
              and profile badge.
            </Typography>
          </Box>

          <Grid container spacing={2} sx={{ mt: 0 }}>
            <Grid item xs={12} sm={6} md={4}>
              <Box
                sx={injectActiveStyle(USER_SUBSCRIPTION_SUPPORTER, {
                  py: 2,
                  px: 3,
                  borderTop: '2px solid #596b95',
                  backgroundImage: 'linear-gradient(#4654755c, #465475)',
                  borderRadius: 2,
                })}>
                {profile.subscription == USER_SUBSCRIPTION_SUPPORTER && <ActiveSubscriptionLabel />}

                <Typography variant="h6">Supporter</Typography>

                <Box component={FeatureList} sx={{ height: 97 }}>
                  <li>Supporter Badge</li>
                  <li>
                    <Link href="#exclusive-features">Refresher Shard</Link>
                  </li>
                </Box>

                {profile.subscription == USER_SUBSCRIPTION_SUPPORTER ? (
                  <ActiveSubscription profile={profile} />
                ) : (
                  <Button
                    variant="outlined"
                    fullWidth
                    sx={{ mt: 10, bgcolor: 'rgb(78, 93, 128)' }}
                    component={Link}
                    href="/subscription/checkout?id=supporter">
                    <Typography variant="h6" sx={{ mr: 0.2 }}>
                      $1
                    </Typography>
                    /mo
                  </Button>
                )}
              </Box>
            </Grid>

            <Grid item xs={12} sm={6} md={4}>
              <Box
                sx={injectActiveStyle(USER_SUBSCRIPTION_TRADER, {
                  py: 2,
                  px: 3,
                  borderTop: '2px solid #629cbd',
                  backgroundImage: 'linear-gradient(#578ba863, #578ba8)',
                  borderRadius: 2,
                })}>
                {profile.subscription == USER_SUBSCRIPTION_TRADER && <ActiveSubscriptionLabel />}

                <Typography variant="h6">Trader</Typography>

                <Box component={FeatureList} sx={{ height: 97 }}>
                  <li>Trader Badge</li>
                  <li>
                    <Link href="#exclusive-features">Refresher Orb</Link>
                  </li>
                </Box>

                {profile.subscription == USER_SUBSCRIPTION_TRADER ? (
                  <ActiveSubscription profile={profile} />
                ) : (
                  <Button
                    variant="outlined"
                    fullWidth
                    sx={{ mt: 10, bgcolor: 'rgb(100, 159, 192)' }}
                    component={Link}
                    href="/subscription/checkout?id=trader">
                    <Typography variant="h6" sx={{ mr: 0.2 }}>
                      $3
                    </Typography>
                    /mo
                  </Button>
                )}
              </Box>
            </Grid>

            <Grid item xs={12} sm={12} md={4}>
              <Box
                sx={injectActiveStyle(USER_SUBSCRIPTION_PARTNER, {
                  py: 2,
                  px: 3,
                  borderTop: '2px solid #ae7f1e',
                  backgroundImage: 'linear-gradient(#a6791d63, #a6791d)',
                  borderRadius: 2,
                  maxWidth: 500,
                  margin: 'auto',
                })}>
                {profile.subscription == USER_SUBSCRIPTION_PARTNER && <ActiveSubscriptionLabel />}

                <Typography variant="h6">Partner</Typography>

                <Box component={FeatureList}>
                  <li>Partner Badge</li>
                  <li>
                    <Link href="#exclusive-features">Refresher Orb</Link>
                  </li>
                  <li>
                    <Link href="#exclusive-features">Shopkeeper&apos;s Contract</Link>
                  </li>
                  <li style={{ textDecoration: 'line-through' }}>
                    <Link href="#">Dedicated Pos-5</Link>
                  </li>
                </Box>

                {profile.subscription == USER_SUBSCRIPTION_PARTNER ? (
                  <ActiveSubscription profile={profile} />
                ) : (
                  <Button
                    variant="outlined"
                    fullWidth
                    sx={{ mt: 10, bgcolor: 'rgb(197, 144, 35)' }}
                    component={Link}
                    href="/subscription/checkout?id=partner">
                    <Typography variant="h6" sx={{ mr: 0.2 }}>
                      $20
                    </Typography>
                    /mo
                  </Button>
                )}
              </Box>
            </Grid>
          </Grid>
          <br />
          <Typography variant="body2" paragraph textAlign="center" color="text.secondary">
            Subscriptions automatically renew and you can cancel your subscription on your
            Paypal&apos;s dashboard.
          </Typography>

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }} id="exclusive-features">
              Exclusive Features
            </Typography>
            <Grid container spacing={1.5}>
              <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 2 }}>
                  <Image
                    src="/assets/refresher-shard.png"
                    alt="assets/refresher-shard.png"
                    width={66}
                    height={48}
                  />
                  <Typography>Refresher Shard</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Automatically refreshes expiring buy orders
                  </Typography>
                </Box>
              </Grid>

              <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 2 }}>
                  <Image
                    src="/assets/refresher-orb.png"
                    alt="assets/refresher-orb.png"
                    width={66}
                    height={48}
                  />
                  <Typography>Refresher Orb</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Automatically refreshes expiring buy orders and listings
                  </Typography>
                </Box>
              </Grid>

              <Grid item md={3} sm={4} xs={6}>
                <Box textAlign="center" sx={{ bgcolor: 'background.paper', p: 2, borderRadius: 2 }}>
                  <Image src="/assets/recipe.png" alt="assets/recipe.png" height={48} width={66} />
                  <Typography>Shopkeeper&apos;s Contract</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Grants the ability to resell items outside your inventory
                  </Typography>
                </Box>
              </Grid>

              <Grid item md={3} sm={4} xs={6}>
                <Box
                  textAlign="center"
                  sx={{
                    bgcolor: 'background.paper',
                    p: 2,
                    borderRadius: 2,
                    filter: 'grayscale(100%)',
                  }}>
                  <Image
                    src="/assets/courier.png"
                    alt="assets/courier.png"
                    height={48}
                    width={66}
                  />
                  <Typography>Dedicated Pos-5</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Exclusive support channel on Discord and Steam
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}

function FeatureUnlockables() {
  return (
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
  )
}

function ActiveSubscription({ profile }: { profile: Profile }) {
  return (
    <Typography variant="subtitle2" sx={{ mt: 12 }}>
      You have active {profile.subscription_type} subscription{' '}
      {profile.subscription_ends_at
        ? `until ${dateCalendar(profile.subscription_ends_at)}`
        : 'automatically billed monthly'}
      .
    </Typography>
  )
}

function ActiveSubscriptionLabel() {
  return (
    <Typography variant="subtitle2" color="yellowgreen" sx={{ float: 'right' }}>
      Active
    </Typography>
  )
}
