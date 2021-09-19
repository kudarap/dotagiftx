import React from 'react'
import PropTypes from 'prop-types'
import { useTheme } from '@material-ui/core/styles'
import useMediaQuery from '@material-ui/core/useMediaQuery'
import * as Auth from '@/service/auth'
import { blacklistSearch } from '@/service/api'
import Theme from '@/components/Theme'
import AppContext from '@/components/AppContext'
// import WhatsNewDialog from '@/components/WhatsNewDialog'
// import SurveyFab from '@/components/SurveyFab'
// import { REPORT_LABEL_SURVEY_NEXT } from '@/constants/report'

function Root({ children }) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('xs'))
  const isTablet = useMediaQuery(theme.breakpoints.down('sm'))
  const currentAuth = Auth.get()
  const isLoggedIn = Auth.isOk()

  const [latestBan, setLatestBan] = React.useState(null)
  React.useEffect(() => {
    ;(async () => {
      try {
        const user = await blacklistSearch({ limit: 1, sort: 'updated_at:desc' })
        if (user) {
          setLatestBan(user[0])
        }
      } catch (error) {
        console.log('error getting lastest ban', error)
      }
    })()
  }, [])

  return (
    <AppContext.Provider value={{ isMobile, isTablet, currentAuth, isLoggedIn, latestBan }}>
      <Theme>{children}</Theme>
      {/* {currentAuth.user_id && (
        <Theme>
          <SurveyFab userID={currentAuth.user_id} label={REPORT_LABEL_SURVEY_NEXT} />
        </Theme>
      )} */}

      {/* {currentAuth.user_id && (
        <Theme>
          <WhatsNewDialog userID={currentAuth.user_id} open />
        </Theme>
      )} */}
    </AppContext.Provider>
  )
}

Root.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Root
