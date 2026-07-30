import React, { useMemo, useState } from 'react'
import PropTypes from 'prop-types'
import { useTheme } from '@mui/material/styles'
import useMediaQuery from '@mui/material/useMediaQuery'
import { PayPalScriptProvider } from '@paypal/react-paypal-js'
import * as Auth from '@/service/auth'
import AppContext from '@/components/AppContext'

const PAYPAL_CLIENT_ID = process.env.NEXT_PUBLIC_PAYPAL_CLIENT_ID

function Root({ children }) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))
  const isTablet = useMediaQuery(theme.breakpoints.down('md'))
  const currentAuth = Auth.get()
  const isLoggedIn = Auth.isOk()

  const [isClient, setIsClient] = useState(false)

  React.useEffect(() => {
    setIsClient(true)
  }, [])

  const contextValue = useMemo(
    () => ({
      isMobile,
      isTablet,
      currentAuth,
      isLoggedIn,
    }),
    [isMobile, isTablet, currentAuth, isLoggedIn]
  )

  // Don't render anything on the server-side to prevent hydration mismatches
  if (!isClient) {
    return null
  }

  return (
    <AppContext.Provider value={contextValue}>
      <PayPalScriptProvider
        options={{
          'client-id': PAYPAL_CLIENT_ID,
          components: 'buttons',
          intent: 'subscription',
          vault: true,
        }}>
        {children}
      </PayPalScriptProvider>
    </AppContext.Provider>
  )
}

Root.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Root
