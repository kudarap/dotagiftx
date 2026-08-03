import React, { useMemo } from 'react'
import PropTypes from 'prop-types'
import { useTheme } from '@mui/material/styles'
import useMediaQuery from '@mui/material/useMediaQuery'
import * as Auth from '@/service/auth'
import AppContext from '@/components/AppContext'

function Root({ children }) {
  const theme = useTheme()
  // useMediaQuery is hydration-safe: it uses defaultMatches (false) on the
  // server and first client render, then updates after mount. Rendering
  // children unconditionally keeps SSR intact.
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))
  const isTablet = useMediaQuery(theme.breakpoints.down('md'))
  const currentAuth = Auth.get()
  const isLoggedIn = Auth.isOk()

  const contextValue = useMemo(
    () => ({
      isMobile,
      isTablet,
      currentAuth,
      isLoggedIn,
    }),
    [isMobile, isTablet, currentAuth, isLoggedIn]
  )

  return <AppContext.Provider value={contextValue}>{children}</AppContext.Provider>
}

Root.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Root
