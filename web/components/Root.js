import React from 'react'
import PropTypes from 'prop-types'
import { useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import * as Auth from '@/service/auth'
import Theme from '@/components/Theme'
import AppContext from '@/components/AppContext'
// import SurveyFab from '@/components/SurveyFab'
// import { REPORT_LABEL_SURVEY_NEXT } from '@/constants/report'

function Root({ children }) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))
  const currentAuth = Auth.get()
  const isLoggedIn = Auth.isOk()

  return (
    <AppContext.Provider value={{ isMobile, currentAuth, isLoggedIn }}>
      <Theme>{children}</Theme>
      {/* {currentAuth.user_id && (
        <Theme>
          <SurveyFab userID={currentAuth.user_id} label={REPORT_LABEL_SURVEY_NEXT} />
        </Theme>
      )} */}
    </AppContext.Provider>
  )
}

Root.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Root
