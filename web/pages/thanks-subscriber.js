import React, { useContext } from 'react'
import { useRouter } from 'next/router'
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
import { APP_CACHE_PROFILE } from '@/constants/app'
import * as Storage from '@/service/storage'
import Avatar from '@/components/Avatar'
import { CDN_URL, processMySubscription, myProfile } from '@/service/api'
import SubscriberBadge from '@/components/SubscriberBadge'

const avatarSize = 92

export default function ThanksSubscriber() {
  const { isLoggedIn } = useContext(AppContext)

  const router = useRouter()
  const subscriptionName = router?.query?.sub
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

  const [verified, setVerified] = React.useState(null)
  React.useEffect(() => {
    if (!subscriptionID) {
      return
    }

    ;(async () => {
      try {
        const res = await processMySubscription(subscriptionID)
        Storage.save(APP_CACHE_PROFILE, res)
        setVerified(true)
      } catch (e) {
        setVerified(false)
      }
    })()
  }, [subscriptionID])

  return (
    <div className="container">
      <Header />

      <main>
        <Container>
          <Box sx={{ mt: 8, textAlign: 'center' }}>
            <Typography variant="h4" component="h1" fontWeight="bold" gutterBottom>
              Thank you for Supporting DotagiftX
            </Typography>
            <Typography>Effect may take few minutes and might ask you to re-login</Typography>
          </Box>

          {profile && subscriptionName && (
            <Box sx={{ mt: 8, textAlign: 'center' }}>
              <Avatar
                sx={{ m: 'auto', mb: 1 }}
                badge={subscriptionName}
                style={{ width: avatarSize, height: avatarSize }}
                src={`${CDN_URL}/${profile.avatar}`}
                // {...retinaSrcSet(profile.avatar, avatarSize, avatarSize)}
              />
              {verified === null && <Typography>Verifying...</Typography>}
              {verified === false && (
                <Typography color="error">Error verifying your subscription</Typography>
              )}
              {verified === true && <SubscriberBadge type={subscriptionName} size="medium" />}
            </Box>
          )}

          <Divider sx={{ my: 4 }} />

          <Box sx={{ textAlign: 'center' }}>
            <Button
              startIcon={<DiscordIcon />}
              variant="outlined"
              size="large"
              component={Link}
              target="_blank"
              rel="noreferrer noopener"
              href="https://discord.gg/UFt9Ny42kM">
              Join our Discord
            </Button>
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
