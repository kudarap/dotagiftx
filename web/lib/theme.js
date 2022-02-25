import { createTheme } from '@mui/material/styles'
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
  },
}

const muiTheme = createTheme(baseThemeOpts)

export default muiTheme

export const muiLightTheme = createTheme({
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
