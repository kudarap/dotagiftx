import { withStyles } from 'tss-react/mui'
import { teal as primary } from '@mui/material/colors'
import Button from '@/components/Button'

export default withStyles(Button, theme => ({
  root: {
    color: theme.palette.getContrastText(primary[900]),
    backgroundColor: primary[800],
    '&:hover': {
      backgroundColor: primary[700],
    },
  },
}))
