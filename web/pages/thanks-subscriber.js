import React from 'react'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Link from '@/components/Link'
import Button from '@/components/Button'
import DiscordIcon from '@/components/DiscordIcon'

export default function ThanksSubscriber() {
  return (
    <div className="container">
      <Header />

      <main>
        <Container>
          <Box sx={{ mt: 8, mb: 4, textAlign: 'center' }}>
            <Typography variant="h4" component="h1" fontWeight="bold" gutterBottom>
              Thank you for Supporting DotagiftX
            </Typography>
            <Typography>Effect may take few minutes and might ask you to re-login</Typography>
          </Box>

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
