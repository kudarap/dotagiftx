import withStyles from '@mui/styles/withStyles'
import { teal as primary } from '@mui/material/colors'
import Button from '@/components/Button'

export default withStyles(theme => ({
  root: {
    color: primary[300],
    borderColor: 'rgba(77, 182, 172, 0.5)',
    // color: theme.palette.getContrastText(primary[900]),
    // backgroundColor: primary[900],
    '&:hover': {
      backgroundColor: 'rgba(77, 182, 172, 0.1)',
    },
  },
}))(Button)
