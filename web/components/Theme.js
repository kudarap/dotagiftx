import React from 'react'
import PropTypes from 'prop-types'

import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider, createTheme, StyledEngineProvider } from '@mui/material/styles'
import { teal, blueGrey } from '@mui/material/colors'

const baseThemeOpts = {
  typography: {
    fontFamily: 'Ubuntu, sans-serif',
  },
  palette: {
    mode: 'dark',
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
  components: {
    MuiAvatar: {
      defaultProps: {
        variant: 'rounded',
      },
    },
  },
}

export const muiTheme = createTheme(baseThemeOpts)

export default function Theme({ children }) {
  return (
    <StyledEngineProvider injectFirst>
      <ThemeProvider theme={muiTheme}>
        <CssBaseline />
        {children}
      </ThemeProvider>
    </StyledEngineProvider>
  )
}
Theme.propTypes = {
  children: PropTypes.node.isRequired,
}

const muiLightTheme = createTheme({
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
