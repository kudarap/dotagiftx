import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider } from '@mui/material/styles'
import { AppCacheProvider } from '@mui/material-nextjs/v15-pagesRouter'
import { CacheProvider } from '@emotion/react'
import { APP_NAME } from '@/constants/strings'
import theme from '@/lib/theme'
import createEmotionCache from '@/lib/createEmotionCache'
import Root from '@/components/Root'
import '@/components/Avatar.css'

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache

export function MyApp2(props) {
  const { Component, pageProps } = props

  return (
    <AppCacheProvider {...props}>
      <Head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <ThemeProvider theme={theme}>
        {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
        <CssBaseline />
        <Component {...pageProps} />
      </ThemeProvider>
    </AppCacheProvider>
  )
}

export default MyApp2

export function MyApp(props) {
  const { Component, emotionCache = clientSideEmotionCache(), pageProps } = props

  return (
    <AppCacheProvider {...props}>
      <CacheProvider value={emotionCache}>
        <Head>
          <meta charSet="UTF-8" />
          <title>{`${APP_NAME} :: Dota 2 Giftables Community Market`}</title>
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
          </Root>
        </ThemeProvider>
      </CacheProvider>
    </AppCacheProvider>
  )
}

MyApp.propTypes = {
  Component: PropTypes.elementType.isRequired,
  emotionCache: PropTypes.object,
  pageProps: PropTypes.object.isRequired,
}
