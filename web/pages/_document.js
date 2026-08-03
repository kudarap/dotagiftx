import * as React from 'react'
import { Html, Head, Main, NextScript } from 'next/document'
import { DocumentHeadTags, documentGetInitialProps } from '@mui/material-nextjs/v15-pagesRouter'
import theme from '@/lib/theme'

const apiOrigin = process.env.NEXT_PUBLIC_API_URL
const cdnOrigin = process.env.NEXT_PUBLIC_CDN_URL

export default function MyDocument(props) {
  return (
    <Html lang="en">
      <Head>
        <DocumentHeadTags {...props} />

        {/* resolves dns for fast load time from other resources */}
        <link rel="preconnect" href="https://cdn.steamstatic.com" />
        {apiOrigin && <link rel="preconnect" href={apiOrigin} />}
        {cdnOrigin && <link rel="preconnect" href={cdnOrigin} />}

        {/* PWA primary color */}
        <meta name="theme-color" content={theme.palette.primary.main} />
        <link rel="icon" href="/favicon.ico" />
        <meta name="emotion-insertion-point" content="" />

        {process.env.NEXT_PUBLIC_UMAMI && (
          <>
            <script
              defer
              src="https://umami.chiligarlic.com/adobo.js"
              data-website-id={process.env.NEXT_PUBLIC_UMAMI}
              data-performance="true"
            />
            <script
              defer
              src="https://umami.chiligarlic.com/recorder.js"
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
