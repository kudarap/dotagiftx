import React, { useContext } from 'react'
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
import { CDN_URL } from '@/service/api'
import SubscriberBadge from '@/components/SubscriberBadge'

const avatarSize = 92

export default function ThanksSubscriber() {
  const { isLoggedIn } = useContext(AppContext)

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

          {profile && (
            <Box sx={{ mt: 8, textAlign: 'center' }}>
              <Avatar
                sx={{ m: 'auto' }}
                badge="supporter"
                style={{ width: avatarSize, height: avatarSize }}
                src={`${CDN_URL}/${profile.avatar}`}
                // {...retinaSrcSet(profile.avatar, avatarSize, avatarSize)}
              />
              <SubscriberBadge type="supporter" size="medium" />
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
