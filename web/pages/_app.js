import * as React from 'react'
import Head from 'next/head'
import { AppCacheProvider } from '@mui/material-nextjs/v15-pagesRouter'
import { ThemeProvider } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import { APP_NAME } from '@/constants/strings'
import theme from '@/lib/theme'
import Root from '@/components/Root'
import '@/components/Avatar.css'

export default function MyApp(props) {
  const { Component, pageProps } = props

  return (
    <AppCacheProvider {...props}>
      <Head>
        <meta charSet="UTF-8" />
        <title>{`${APP_NAME} :: Dota 2 Giftables Community Market`}</title>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=6.0" />
      </Head>
      <ThemeProvider theme={theme}>
        {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
        <CssBaseline />

        <Root>
          <Component {...pageProps} />
        </Root>
      </ThemeProvider>
    </AppCacheProvider>
  )
}
