import React from 'react'
import Document, { Head, Html, Main, NextScript } from 'next/document'
import { withEmotionCache } from 'tss-react/nextJs'
import muiTheme from '@/lib/theme'
import createEmotionCache from '@/lib/createEmotionCache'

class MyDocument extends Document {
  render() {
    return (
      <Html lang="en">
        <Head>
          {/* resolves dns for fast load time from other resources */}
          <link rel="preconnect" href="https://fonts.googleapis.com" />
          <link rel="preconnect" href="https://www.googleanalytics.com" />
          <link rel="preconnect" href="https://www.googletagmanager.com" />
          <link rel="preconnect" href="https://cdn.cloudflare.steamstatic.com" />

          {/* PWA primary color */}
          <meta name="theme-color" content={muiTheme.palette.primary.main} />
          <link
            href="https://fonts.googleapis.com/css2?family=Ubuntu:wght@300;400;500;700&display=swap"
            rel="stylesheet"
          />
          <link rel="icon" href="/favicon.ico" />

          {/* Inject MUI styles first to match with the prepend: true configuration. */}
          {this.props.emotionStyleTags}

          {process.env.NEXT_PUBLIC_GA && (
            <>
              <script
                async
                src={`https://www.googletagmanager.com/gtag/js?id=${process.env.NEXT_PUBLIC_GA}`}
              />
              <script
                dangerouslySetInnerHTML={{
                  __html: `
                window.dataLayer = window.dataLayer || [];
                function gtag(){window.dataLayer.push(arguments)}
                gtag("js", new Date());
                gtag("config", "${process.env.NEXT_PUBLIC_GA}");`,
                }}
              />
            </>
          )}
        </Head>
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    )
  }
}

export default withEmotionCache({
  Document: MyDocument,
  // Every emotion cache used in the app should be provided.
  // Caches for MUI should use "prepend": true.
  getCaches: () => [createEmotionCache()],
})
