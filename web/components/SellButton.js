import { withStyles } from '@material-ui/core/styles'
import primary from '@material-ui/core/colors/teal'
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
