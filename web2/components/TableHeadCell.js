import { withStyles } from 'tss-react/mui'
import TableCell from '@mui/material/TableCell'

export default withStyles(TableCell, theme => ({
  head: {
    textTransform: 'uppercase',
    color: theme.palette.text.secondary,
  },
}))
