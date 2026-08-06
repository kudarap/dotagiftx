import PropTypes from 'prop-types'
import { withStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'

const Code = props => {
  return <Typography {...props} />
}
Code.propTypes = {
  component: PropTypes.string,
}
Code.defaultProps = {
  component: 'span',
}

export default withStyles(Code, () => ({
  root: {
    fontFamily: 'monospace',
  },
}))
