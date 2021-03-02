import React from 'react'
import PropTypes from 'prop-types'
import { useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import * as Auth from '@/service/auth'
import Theme from '@/components/Theme'
import AppContext from '@/components/AppContext'
import { useRouter } from 'next/router'

function Root({ children }) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))
  const currentAuth = Auth.get()
  const isLoggedIn = Auth.isOk()

  // Checks for new auth field steam_id.
  // need to re-authenticate to get new field.
  const router = useRouter()
  React.useEffect(() => {
    if (!isLoggedIn || currentAuth.steam_id) {
      return
    }

    Auth.clear()
    router.push('/login')
  }, [currentAuth])

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
