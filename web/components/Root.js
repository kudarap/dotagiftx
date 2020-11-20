import React from 'react'
import PropTypes from 'prop-types'
import { useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import { get, isOk } from '@/service/auth'
import Theme from '@/components/Theme'
import AppContext from '@/components/AppContext'

function Root({ children }) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))
  const currentAuth = get()
  const isLoggedIn = isOk()
  return (
    <AppContext.Provider value={{ isMobile, currentAuth, isLoggedIn }}>
      <Theme>{children}</Theme>
    </AppContext.Provider>
  )
}

Root.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Root
