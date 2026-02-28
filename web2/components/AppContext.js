import { createContext } from 'react'

const AppContext = createContext({
  currentAuth: null,
  latestBan: null,
  isLoggedIn: false,
  isMobile: false,
  isTablet: false,
})

export default AppContext
