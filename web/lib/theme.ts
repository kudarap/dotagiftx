import { Ubuntu } from 'next/font/google'
import { createTheme, type PaletteOptions, type ThemeOptions } from '@mui/material/styles'
import { teal, blueGrey, grey } from '@mui/material/colors'
import { responsiveFontSizes } from '@mui/material'

declare module '@mui/material/styles' {
  interface Palette {
    bid: Palette['primary']
    accent: Palette['primary']
    app: { white: string }
  }
  interface PaletteOptions {
    bid?: PaletteOptions['primary']
    accent?: PaletteOptions['primary']
    app?: { white: string }
  }
}

declare module '@mui/material/Button' {
  interface ButtonPropsVariantOverrides {
    default: true
  }
}

declare module '@mui/material/Button' {
  interface ButtonPropsColorOverrides {
    bid: true
    accent: true
  }
}

const font = Ubuntu({
  weight: ['300', '400', '500', '700'],
  subsets: ['latin'],
  display: 'swap',
})

const defaultPalette: PaletteOptions = {
  mode: 'dark',
  primary: {
    main: grey[200],
    light: grey[100],
    dark: grey[400],
  },
  secondary: {
    main: '#C79123',
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
}

const darkCarnivalPalette: PaletteOptions = {
  secondary: {
    main: '#FFDFB7',
  },
  background: {
    default: '#090a19',
    paper: '#1e2331ef',
  },
}

const baseThemeOpts: ThemeOptions = {
  typography: {
    fontFamily: font.style.fontFamily,
  },
  palette: {
    ...defaultPalette,
    ...darkCarnivalPalette,
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
    MuiTableContainer: {
      styleOverrides: {
        root: {
          border: '1px solid #515051',
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        root: {
          // borderBottomColor: '#263238',
          borderBottomColor: '#51505161',
        },
      },
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
