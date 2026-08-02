import React, { useContext } from 'react'
import { useRouter } from 'next/router'
import { APP_CACHE_PROFILE } from '@/constants/app'
import { USER_SUBSCRIPTION_MAP_LABEL } from '@/constants/user'
import * as Storage from '@/service/storage'
import { CDN_URL, processMySubscription, myProfile } from '@/service/api'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import Divider from '@mui/material/Divider'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DiscordIcon from '@/components/DiscordIcon'
import AppContext from '@/components/AppContext'
import Avatar from '@/components/Avatar'
import SubscriberBadge from '@/components/SubscriberBadge'

const avatarSize = 184 * 0.75

const backlightColor = {
  partner: '#ffab00',
  trader: '#06a5ff',
  supporter: '#034eff',
}

export default function ThanksSubscriber() {
  const { isLoggedIn } = useContext(AppContext)

  const router = useRouter()
  const subscriptionID = router?.query?.subid

  const [profile, setProfile] = React.useState(null)
  React.useEffect(() => {
    ;(async () => {
      if (!isLoggedIn) {
        return
      }

      const cached = Storage.get(APP_CACHE_PROFILE)
      if (cached) {
        setProfile(cached)
        return
      }

      const remote = await myProfile.GET()
      Storage.save(APP_CACHE_PROFILE, remote)
      setProfile(remote)
    })()
  }, [])

  const [subscription, setSubscription] = React.useState(null)
  React.useEffect(() => {
    if (!subscriptionID) {
      return
    }

    ;(async () => {
      try {
        const res = await processMySubscription(subscriptionID)
        Storage.save(APP_CACHE_PROFILE, res)
        setSubscription(String(USER_SUBSCRIPTION_MAP_LABEL[res.subscription]).toLowerCase())
      } catch (e) {
        setSubscription(false)
      }
    })()
  }, [subscriptionID])

  return (
    <div className="container">
      <Header />

      <main>
        <Container>
          <Box sx={{ mt: 8, textAlign: 'center', visibility: subscription ? 'inherit' : 'hidden' }}>
            <Typography
              sx={{ mt: 8 }}
              style={{
                background: 'linear-gradient( to right, #CB8F37 20%, #F0CF59 50%, #B5793D 80% )',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                letterSpacing: 3,
              }}
              variant="h3"
              component="h1"
              fontWeight="bold"
              color="secondary">
              Thank you!
            </Typography>
            <Typography
              style={{
                background: 'linear-gradient( to right, #CB8F37 20%, #F0CF59 50%, #B5793D 80% )',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                letterSpacing: 3,
              }}
              variant="h6">
              for keeping the servers on
            </Typography>
          </Box>

          {profile && subscriptionID && (
            <Box sx={{ mt: 8, textAlign: 'center' }}>
              <Avatar
                sx={{ m: 'auto', mb: 1 }}
                badge={subscription}
                style={{ width: avatarSize, height: avatarSize }}
                src={`${CDN_URL}/${profile.avatar}`}
              />
              {subscription === null && <Typography>Verifying...</Typography>}
              {subscription === false && (
                <>
                  <Typography color="error" sx={{ mb: 1 }}>
                    Error verifying your subscription
                  </Typography>
                  <Button
                    variant="outlined"
                    size="large"
                    onClick={() => {
                      router.reload()
                    }}>
                    Try again
                  </Button>
                </>
              )}
              {subscription && <SubscriberBadge type={subscription} size="large" />}

              <br />
              <br />
            </Box>
          )}

          <Divider
            sx={{
              mt: 20,
              mb: 4,
              transition: 'all 2s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
              boxShadow: subscription
                ? `0 -60px 290px ${backlightColor[subscription]}`
                : `0 -60px 290px black`,
            }}
          />

          <Box sx={{ textAlign: 'center', display: subscription ? 'inherit' : 'none' }}>
            <Button
              sx={{ mr: 2 }}
              startIcon={<DiscordIcon />}
              variant="outlined"
              size="large"
              component={Link}
              target="_blank"
              rel="noreferrer noopener"
              href="https://discord.gg/UFt9Ny42kM">
              Join our Discord
            </Button>
            <Button
              color="secondary"
              variant="outlined"
              size="large"
              component={Link}
              href={`/profiles/${profile?.steam_id}`}>
              Check your profile
            </Button>
          </Box>

          {subscription && (
            <Typography align="center" sx={{ mt: 4 }}>
              You&apos;re subscription is now active, badge may take a few minutes to reflect on
              your profile.
              <br />
              Please reach out to discord if there&apos;s any issue, Enjoy!
            </Typography>
          )}
        </Container>
      </main>

      <Footer />
    </div>
  )
}
