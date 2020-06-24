import { withStyles } from '@material-ui/core/styles'
import green from '@material-ui/core/colors/lightGreen'
import Button from '@/components/Button'

const BuyButton = withStyles(theme => ({
  root: {
    color: theme.palette.getContrastText(green[700]),
    backgroundColor: green[700],
    '&:hover': {
      backgroundColor: green[800],
    },
  },
}))(Button)

export default BuyButton
