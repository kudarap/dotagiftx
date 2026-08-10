import * as React from 'react'
import { Head, Html, Main, NextScript } from 'next/document'
import { documentGetInitialProps, DocumentHeadTags } from '@mui/material-nextjs/v15-pagesRouter'
import theme from '@/lib/theme'

export default function MyDocument(props) {
  return (
    <Html lang="en">
      <Head>
        <DocumentHeadTags {...props} />

        {/* resolves dns for fast load time from other resources */}
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://www.googleanalytics.com" />
        <link rel="preconnect" href="https://www.googletagmanager.com" />
        <link rel="preconnect" href="https://cdn.steamstatic.com" />

        {/* PWA primary color */}
        <meta name="theme-color" content={theme.palette.primary.main} />
        <link rel="icon" href="/favicon.ico" />
        <meta name="emotion-insertion-point" content="" />

        {process.env.NEXT_PUBLIC_UMAMI_URL && process.env.NEXT_PUBLIC_UMAMI && (
          <>
            <script
              defer
              src={`${process.env.NEXT_PUBLIC_UMAMI_URL}/script.js`}
              data-website-id={process.env.NEXT_PUBLIC_UMAMI}
              data-performance="true"
            />
            <script
              defer
              src={`${process.env.NEXT_PUBLIC_UMAMI_URL}/recorder.js`}
              data-website-id={process.env.NEXT_PUBLIC_UMAMI}
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

MyDocument.getInitialProps = async ctx => {
  const finalProps = await documentGetInitialProps(ctx)
  return finalProps
}
