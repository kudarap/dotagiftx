import withStyles from '@mui/styles/withStyles'
import primary from '@mui/material/colors/lightGreen'
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
