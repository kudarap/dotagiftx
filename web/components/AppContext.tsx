import { createContext } from 'react'
import type { Auth } from '@/service/auth'

export interface AppContextValue {
  currentAuth: Auth | null
  latestBan: string | null
  isLoggedIn: boolean
  isMobile: boolean
  isTablet: boolean
}

const AppContext = createContext<AppContextValue>({
  currentAuth: null,
  latestBan: null,
  isLoggedIn: false,
  isMobile: false,
  isTablet: false,
})

export default AppContext
