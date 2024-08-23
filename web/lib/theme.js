import { createTheme } from '@mui/material/styles'
import { teal, blueGrey, grey } from '@mui/material/colors'
import { responsiveFontSizes } from '@mui/material'

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
      main: '#7fbc8b',
    },
    bid: {
      main: teal[300],
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
    MuiAppBar: {
      styleOverrides: {
        root: {
          borderTop: 'none',
          borderLeft: 'none',
          borderRight: 'none',
        },
      },
    },
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
      styleOverrides: {
        root: {
          textTransform: 'none',
        },
      },
      defaultProps: {
        // variant: 'default',
      },
      variants: [
        {
          props: { variant: 'default' },
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
    MuiAlert: {
      variants: [
        {
          props: { variant: 'filled' },
          style: {
            color: 'white',
          },
        },
      ],
    },
  },
}

const muiTheme = createTheme(baseThemeOpts)

export default responsiveFontSizes(muiTheme)

export const muiLightTheme = responsiveFontSizes(
  createTheme({
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
)
