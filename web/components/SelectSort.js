import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import MenuItem from '@material-ui/core/MenuItem'
import Select from '@material-ui/core/Select'
import FormControl from '@material-ui/core/FormControl'

const StyledSelect = withStyles(theme => ({
  root: {
    fontSize: theme.typography.fontSize,
  },
}))(props => <Select {...props} />)

export default function SelectSort({ options, variant, size, ...other }) {
  return (
    <FormControl {...{ variant, size }}>
      <StyledSelect id="select-sort" {...other}>
        {options.map(opt => (
          <MenuItem key={opt.value} value={opt.value}>
            {opt.label}
          </MenuItem>
        ))}
      </StyledSelect>
    </FormControl>
  )
}
SelectSort.propTypes = {
  options: PropTypes.arrayOf(PropTypes.object),
  variant: PropTypes.string,
  size: PropTypes.string,
}
SelectSort.defaultProps = {
  options: [],
  variant: null,
  size: null,
}
