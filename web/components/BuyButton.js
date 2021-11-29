import withStyles from '@mui/styles/withStyles'
import { lightGreen as primary } from '@mui/material/colors'
import Button from '@/components/Button'

export default withStyles(theme => ({
  root: {
    color: theme.palette.getContrastText(primary[900]),
    backgroundColor: primary[900],
    '&:hover': {
      backgroundColor: primary[800],
    },
  },
}))(Button)
