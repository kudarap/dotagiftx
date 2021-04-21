import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

// border: 1px solid goldenrod,
// padding: 0 5px,
// border-radius: 4px,
// color: white,
// background: goldenrod,
// font-weight: 500,

export default withStyles(theme => ({
  root: {
    border: '1px solid goldenrod',
    padding: '0 5px',
    borderRadius: '4px',
    color: 'white',
    background: 'goldenrod',
    display: 'inline',
    fontSize: '0.85rem',
    fontWeight: theme.typography.fontWeightMedium,
    verticalAlign: 'middle',
  },
}))(Typography)
