import React from 'react'
import PropTypes from 'prop-types'
import CssBaseline from '@material-ui/core/CssBaseline'
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles'
import teal from '@material-ui/core/colors/teal'
import { blueGrey } from '@material-ui/core/colors'

const baseThemeOpts = {
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
    accent: {
      main: teal.A200,
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
  overrides: {
    MuiAvatar: {
      root: {
        borderRadius: '15%',
      },
    },
  },
}

export const muiTheme = createMuiTheme(baseThemeOpts)

export default function Theme({ children }) {
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

const muiLightTheme = createMuiTheme({
  ...baseThemeOpts,
  palette: {
    ...baseThemeOpts.palette,
    type: 'light',
    background: {
      paper: blueGrey.A100,
    },
  },
  overrides: {
    MuiTableCell: {
      root: {
        borderBottomColor: blueGrey[200],
      },
    },
  },
})

export function LightTheme({ children }) {
  return <ThemeProvider theme={muiLightTheme}>{children}</ThemeProvider>
}
LightTheme.propTypes = {
  children: PropTypes.node.isRequired,
}
