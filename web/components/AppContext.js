import { createContext } from 'react'

const AppContext = createContext({
  currentAuth: null,
  isLoggedIn: false,
  isMobile: false,
})

export default AppContext
