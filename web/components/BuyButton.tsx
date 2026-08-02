import { withStyles } from 'tss-react/mui'
import { lightGreen as primary } from '@mui/material/colors'
import Button from '@/components/Button'

export default withStyles(Button, theme => ({
  root: {
    color: theme.palette.getContrastText(primary[900]),
    backgroundColor: primary[900],
    '&:hover': {
      backgroundColor: primary[800],
    },
  },
}))
