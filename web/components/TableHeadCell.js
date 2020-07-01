import { withStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'

export default withStyles(theme => ({
  head: {
    textTransform: 'uppercase',
    color: theme.palette.text.secondary,
  },
}))(TableCell)
