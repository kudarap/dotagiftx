import React from 'react'
import Head from 'next/head'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import { APP_NAME } from '@/constants/strings'
import Link from '@/components/Link'

export default function ThanksSubscriber() {
  return (
    <div className="container">
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Data Loss Incident</title>
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
              Data Loss Incident
            </Typography>
            <Typography>Oct 24, 2025</Typography>
            <br />

            <Typography>
              Unfortunately, half of records on DotagiftX has been accidentally deleted and it could
              not be restored. Data loss affects all users and their records on the following:
              <ul>
                <li>Listing</li>
                <li>Reservations</li>
                <li>Orders</li>
                <li>Stats</li>
                <li>History</li>
              </ul>
            </Typography>

            <Typography>
              For the past few weeks, server crash has been frequent due to system resources
              capacity and decided to scale our system and provision for new features in the future
              and the incident happened during the maintenance on Oct 22, 2025.
            </Typography>

            <Typography>
              User authentication was also affected, but already deployed a hotfix to enabled login.
              In any case you still have issue please reach out to{' '}
              <u>
                <Link
                  href="https://discord.gg/JbAm39ubSr"
                  target="_blank"
                  rel="noreferrer noopener">
                  #incident-data-loss-login
                </Link>
              </u>{' '}
              discord channel.
            </Typography>
            <br />

            <Typography>
              I apologize for this grave error and will be more careful next time and make sure
              contingency plan are in place to prevent this major incident from happening. Thank you
              for understanding, if you have more concern you can join on{' '}
              <u>
                <Link
                  href="https://discord.gg/zdwA3zD5NH"
                  target="_blank"
                  rel="noreferrer noopener">
                  #incident-data-loss-concern
                </Link>
              </u>{' '}
              discord channel.
            </Typography>
            <br />

            <Typography>- kudarap</Typography>
          </Box>
        </Container>
      </main>

      <Footer />
    </div>
  )
}
