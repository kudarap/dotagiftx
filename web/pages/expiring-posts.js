import React from 'react'
import Head from 'next/head'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { APP_NAME } from '@/constants/strings'
import { Divider } from '@mui/material'
import Link from '@/components/Link'

export default function ThanksSubscriber() {
  return (
    <div className="container">
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Expiring posts</title>
      </Head>

      <Header />

      <main>
        <Container>
          <Box
            sx={{
              mt: {
                xs: 4,
                md: 8,
              },
              mb: 4,
            }}>
            <Typography variant="h4" component="h1" fontWeight="bold" gutterBottom>
              Expiring posts
            </Typography>
            <Typography>
              Due to the number of outdated listings on the site, some improvements were made to
              make sure that posts were updated regularly by both seller and buyer. We will roll out
              the update on <span style={{ color: 'salmon' }}>May 1, 2022</span> and start removing
              aged posts.
            </Typography>
          </Box>

          <Box sx={{ mb: 2 }}>
            <Typography fontWeight="bold">Listings - 30 days</Typography>
            <Typography color="text.secondary">
              Seller listings will be available on site for 30 days from posted date, they need to
              re-list their items and make necessary updates if needed and will reset the
              expiration.
            </Typography>
            <br />

            <Typography fontWeight="bold">Buy orders - 7 days</Typography>
            <Typography color="text.secondary">
              Buyer orders on the other hand will have 7 days validity and same with the item
              listings, they need to re-post their order if they still require it.
            </Typography>
            <br />
          </Box>

          <Typography>
            Let us know what you think by reaching us on{' '}
            <Link
              underline="always"
              target="_blank"
              rel="noreferrer noopener"
              href="https://discord.gg/UFt9Ny42kM">
              Discord
            </Link>{' '}
            or send us a{' '}
            <Link underline="always" href="/feedback">
              Feedback
            </Link>
            .
          </Typography>
          <br />

          <Typography>
            Also in this update, we added a new{' '}
            <Link href="/treasures" underline="always">
              Treasures
            </Link>{' '}
            page and a subscription-based feature{' '}
            <Link href="/plus" underline="always">
              Dotagift+
            </Link>{' '}
            to support the project.
          </Typography>

          <Divider sx={{ mt: 5, mb: 3.5 }} />

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Frequently Asked Questions
            </Typography>
            <Typography fontWeight="bold">What will happen with my current listings?</Typography>
            <Typography color="text.secondary">
              No change but all your listing older than 30 days on effective date will be removed.
            </Typography>
            <br />

            <Typography fontWeight="bold">What will happen with my current buy orders?</Typography>
            <Typography color="text.secondary">
              No change but all your orders older than 7 days on effective date will be removed.
            </Typography>
            <br />

            <Typography fontWeight="bold">Do I need to re-post after it expires?</Typography>
            <Typography color="text.secondary">
              Yes, as it will be automatically removed after the expiration date.
            </Typography>
            <br />

            <Typography fontWeight="bold">
              Does updating the price resets the expiration?
            </Typography>
            <Typography color="text.secondary">
              Yes, because that is the point of this update. Updating the listing with same pricing
              will also reset the expiration.
            </Typography>
            <br />

            <Typography fontWeight="bold">What should I do to avoid re-posting?</Typography>
            <Typography color="text.secondary">
              We do have new features on site where you need to subscribe a month for your items to
              not expire. Please refer to{' '}
              <Link href="/plus" underline="always">
                Dotagift+
              </Link>{' '}
              page.
            </Typography>
            <br />

            <Typography fontWeight="bold">Why the sudden change?</Typography>
            <Typography color="text.secondary">
              This is to ensure that all items posted on site is updated and active users will be
              prioritize.
            </Typography>
            <br />
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
