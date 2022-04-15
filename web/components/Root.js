import React from 'react'
import PropTypes from 'prop-types'
import { useTheme } from '@mui/material/styles'
import useMediaQuery from '@mui/material/useMediaQuery'
import * as Auth from '@/service/auth'
import { blacklistSearch } from '@/service/api'
import AppContext from '@/components/AppContext'
import WhatsNewDialog from '@/components/WhatsNewDialog'
import { PayPalScriptProvider } from '@paypal/react-paypal-js'
// import SurveyFab from '@/components/SurveyFab'
// import { REPORT_LABEL_SURVEY_NEXT } from '@/constants/report'

const PAYPAL_CLIENT_ID = process.env.NEXT_PUBLIC_PAYPAL_CLIENT_ID

function Root({ children }) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))
  const isTablet = useMediaQuery(theme.breakpoints.down('md'))
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
      <PayPalScriptProvider
        options={{
          'client-id': PAYPAL_CLIENT_ID,
          components: 'buttons',
          intent: 'subscription',
          vault: true,
        }}>
        {children}
      </PayPalScriptProvider>

      {/* {currentAuth.user_id && (
        <Theme>
          <SurveyFab userID={currentAuth.user_id} label={REPORT_LABEL_SURVEY_NEXT} />
        </Theme>
      )} */}

      {/* {currentAuth.user_id && ( */}
      <WhatsNewDialog userID={currentAuth.user_id} />
      {/* )} */}
    </AppContext.Provider>
  )
}

Root.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Root
