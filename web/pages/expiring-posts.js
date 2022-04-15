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
              Due to number of outdated listings on site, some improvements were made to make sure
              that posts were updated regularly by both seller and buyer. We will role out the
              update on <span style={{ color: 'coral' }}>May 1, 2022</span> and start removing very
              old posts.
            </Typography>
          </Box>

          <Box sx={{ mb: 2 }}>
            <Typography fontWeight="bold">Listings - 30 days</Typography>
            <Typography color="text.secondary">
              Sellers item will be available on site for 30 days from posted date, they need to
              re-list their items and make necessary updates if needed.
            </Typography>
            <br />

            <Typography fontWeight="bold">Buy orders - 7 days</Typography>
            <Typography color="text.secondary">
              on the other hand will be having 7 days validity and same with the item listings, they
              need to re-post their order if they still require it.
            </Typography>
            <br />
          </Box>

          <Typography>
            Let us know what you think by reaching us on{' '}
            <Link
              color="secondary"
              target="_blank"
              rel="noreferrer noopener"
              href="https://discord.gg/UFt9Ny42kM">
              Discord
            </Link>{' '}
            or send us a{' '}
            <Link color="secondary" href="/feedback">
              Feedback
            </Link>
            .
          </Typography>
          <br />

          <Typography>
            Also in this update, we added a new <Link href="/treasures">Treasures</Link> page and a
            subscription-based feature <Link href="/plus">Dotagift Plus</Link> to support the
            project.
          </Typography>

          <Divider sx={{ my: 4 }} />

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Frequently Asked Questions
            </Typography>
            <Typography fontWeight="bold">What will happen with my current listings?</Typography>
            <Typography color="text.secondary">
              Your current listing will still be available and 30days(seller)/ 7 days(buyer)
              validity will start on effective date provided for expiring items.
            </Typography>
            <br />

            <Typography fontWeight="bold">
              Do I need to re post my listings/orders after it expires?
            </Typography>
            <Typography color="text.secondary">
              Yes, as it will be automatically deleted after 30days(seller)/ 7 days(buyer).
            </Typography>
            <br />

            <Typography fontWeight="bold">
              What should I do to avoid re-listing/re-posting?
            </Typography>
            <Typography color="text.secondary">
              We do have new features on site where you need to subscribe a month for your items to
              not expire. Please refer to Dotagift+ page.
            </Typography>
            <br />

            <Typography fontWeight="bold">Why the sudden change/update?</Typography>
            <Typography color="text.secondary">
              This is to ensure that all items posted on site is updated and active seller/buyer
              will be prioritize.
            </Typography>
            <br />
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
