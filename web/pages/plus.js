import React from 'react'
import Head from 'next/head'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { APP_NAME } from '@/constants/strings'

function PricingCard() {
  return (
    <Box>
      <Typography variant="h6">Supporter</Typography>
      <Box component="ul" sx={{ height: 98 }}>
        <li>Supporter Badge</li>
        <li>Refresher Shard</li>
      </Box>
      <Button>$1/mo</Button>
    </Box>
  )
}

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

          <Grid container spacing={4} sx={{ mt: 4 }}>
            <Grid item xs={12} sm={6} md={4}>
              <Box
                sx={{
                  p: 4,
                  borderTop: '2px solid rgb(122, 112, 150)',
                  backgroundImage: 'linear-gradient(rgba(79, 73, 96, 0.314), rgb(79, 73, 96))',
                }}>
                <Typography variant="h6">Supporter</Typography>
                <Box component="ul" sx={{ height: 97 }}>
                  <li>Profile Badge</li>
                  <li>Refresher Shard</li>
                </Box>
                <Button variant="outlined" fullWidth>
                  $1/mo
                </Button>
              </Box>
            </Grid>

            <Grid item xs={12} sm={6} md={4}>
              <Box
                sx={{
                  p: 4,
                  borderTop: '2px solid rgb(60, 106, 251)',
                  backgroundImage: 'linear-gradient(rgba(45, 78, 180, 0.314), rgb(45, 78, 180))',
                }}>
                <Typography variant="h6">Trader</Typography>
                <Box component="ul" sx={{ height: 97 }}>
                  <li>Profile Badge</li>
                  <li>Refresher Orb</li>
                </Box>
                <Button variant="outlined" fullWidth>
                  $3/mo
                </Button>
              </Box>
            </Grid>

            <Grid item xs={12} sm={12} md={4}>
              <Box
                sx={{
                  p: 4,
                  maxWidth: 500,
                  margin: 'auto',
                  borderTop: '2px solid rgb(126, 87, 255)',
                  backgroundImage: 'linear-gradient(rgba(87, 65, 162, 0.314), rgb(87, 65, 162))',
                }}>
                <Typography variant="h6">Partner</Typography>
                <Box component="ul">
                  <li>Profile Badge</li>
                  <li>Refresher Orb</li>
                  <li>Shopkeeper's Contract</li>
                  <li>Dedicated Pos-5</li>
                </Box>
                <Button variant="outlined" fullWidth>
                  $20/mo
                </Button>
              </Box>
            </Grid>
          </Grid>
          <br />
          <Typography variant="caption">Subscriptions automatically renew</Typography>

          <Box sx={{ mt: 4 }}>
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
        </Container>
      </main>

      <Footer />
    </div>
  )
}
