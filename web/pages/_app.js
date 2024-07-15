import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider } from '@mui/material/styles'
import { CacheProvider } from '@emotion/react'
import { APP_NAME } from '@/constants/strings'
import theme from '@/lib/theme'
import createEmotionCache from '@/lib/createEmotionCache'
import Root from '@/components/Root'
import '@/components/Avatar.css'
import { Analytics } from '@vercel/analytics/react'

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache

export default function MyApp(props) {
  const { Component, emotionCache = clientSideEmotionCache(), pageProps } = props

  return (
    <CacheProvider value={emotionCache}>
      <Head>
        <meta charset="UTF-8" />
        <title>{APP_NAME} :: Dota 2 Giftables Community Market</title>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=6.0" />
        {/* <meta */}
        {/*  name="viewport" */}
        {/*  content="width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no" */}
        {/* /> */}
      </Head>

      <ThemeProvider theme={theme}>
        <CssBaseline />

        <Root>
          <Component {...pageProps} />
          <Analytics />
        </Root>
      </ThemeProvider>
    </CacheProvider>
  )
}

MyApp.propTypes = {
  Component: PropTypes.elementType.isRequired,
  emotionCache: PropTypes.object,
  pageProps: PropTypes.object.isRequired,
}
