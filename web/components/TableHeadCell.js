import withStyles from '@mui/styles/withStyles'
import TableCell from '@mui/material/TableCell'

export default withStyles(theme => ({
  head: {
    textTransform: 'uppercase',
    color: theme.palette.text.secondary,
  },
}))(TableCell)
