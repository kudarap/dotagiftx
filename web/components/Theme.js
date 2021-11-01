import React from 'react'
import PropTypes from 'prop-types'

import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider, createTheme, StyledEngineProvider } from '@mui/material/styles'
import { teal, blueGrey, grey } from '@mui/material/colors'

const baseThemeOpts = {
  typography: {
    fontFamily: 'Ubuntu, sans-serif',
  },
  palette: {
    mode: 'dark',
    primary: {
      main: grey[200],
      light: grey[100],
      dark: grey[400],
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
    MuiSelect: {
      defaultProps: {
        variant: 'standard',
      },
    },
    MuiButton: {
      defaultProps: {
        // variant: 'default',
      },
      variants: [
        {
          props: { variant: 'defaultx' },
          style: {
            textTransform: 'none',
            border: `2px dashed white`,
          },
        },
      ],
    },
    MuiLink: {
      defaultProps: {
        underline: 'hover',
      },
    },
    MuiPaper: {
      styleOverrides: { root: { backgroundImage: 'unset' } },
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
    mode: 'light',
    background: {
      paper: blueGrey.A100,
    },
  },
  components: {
    MuiTableCell: {
      styleOverrides: {
        root: {
          borderBottomColor: blueGrey[200],
        },
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
