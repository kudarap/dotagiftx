import React from 'react'
import PropTypes from 'prop-types'
import Head from 'next/head'
import { APP_NAME } from '@/constants/strings'
import Root from '@/components/Root'
import '@/components/Avatar.css'

export default function MyApp(props) {
  const { Component, pageProps } = props

  React.useEffect(() => {
    // Remove the server-side injected CSS.
    // eslint-disable-next-line no-undef
    const jssStyles = document.querySelector('#jss-server-side')
    if (jssStyles) {
      jssStyles.parentElement.removeChild(jssStyles)
    }
  }, [])

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Dota 2 Giftables Community Market</title>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=6.0" />
        {/* <meta */}
        {/*  name="viewport" */}
        {/*  content="width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no" */}
        {/* /> */}
      </Head>
      <Root>
        <Component {...pageProps} />
      </Root>
    </>
  )
}

MyApp.propTypes = {
  Component: PropTypes.elementType.isRequired,
  pageProps: PropTypes.object.isRequired,
}
