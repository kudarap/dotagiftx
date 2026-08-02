import React, { useMemo, useState } from 'react'
import { useTheme } from '@mui/material/styles'
import useMediaQuery from '@mui/material/useMediaQuery'
import * as Auth from '@/service/auth'
import AppContext from '@/components/AppContext'
import type { AppContextValue } from '@/components/AppContext'

function Root({ children }: { children: React.ReactNode }) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))
  const isTablet = useMediaQuery(theme.breakpoints.down('md'))
  const currentAuth = Auth.get()
  const isLoggedIn = Auth.isOk()

  const [isClient, setIsClient] = useState(false)

  React.useEffect(() => {
    setIsClient(true)
  }, [])

  const contextValue = useMemo<AppContextValue>(
    () => ({
      isMobile,
      isTablet,
      currentAuth,
      isLoggedIn,
      latestBan: null,
    }),
    [isMobile, isTablet, currentAuth, isLoggedIn]
  )

  // Don't render anything on the server-side to prevent hydration mismatches
  if (!isClient) {
    return null
  }

  return <AppContext.Provider value={contextValue}>{children}</AppContext.Provider>
}

export default Root
