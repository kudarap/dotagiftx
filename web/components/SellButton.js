import withStyles from '@mui/styles/withStyles'
import { teal as primary } from '@mui/material/colors'
import Button from '@/components/Button'

export default withStyles(theme => ({
  root: {
    color: theme.palette.getContrastText(primary[900]),
    backgroundColor: primary[800],
    '&:hover': {
      backgroundColor: primary[700],
    },
  },
}))(Button)
