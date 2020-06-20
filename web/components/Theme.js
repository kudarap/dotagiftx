import React from 'react'
import PropTypes from 'prop-types'
import CssBaseline from '@material-ui/core/CssBaseline'
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles'

export const muiTheme = createMuiTheme({
  typography: {
    fontFamily: 'Ubuntu, sans-serif',
  },
  palette: {
    type: 'dark',
    primary: {
      main: '#263238',
    },
    secondary: {
      main: '#C79123',
    },
    background: {
      default: '#263238',
      paper: '#2e3d44',
    },
    // App specific colors.
    app: {
      white: '#FFFBF1',
    },
  },
})

function Theme({ children }) {
  return (
    <ThemeProvider theme={muiTheme}>
      <CssBaseline />
      {children}
    </ThemeProvider>
  )
}
Theme.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Theme
