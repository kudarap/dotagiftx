import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import MenuItem from '@material-ui/core/MenuItem'
import Select from '@material-ui/core/Select'
import FormControl from '@material-ui/core/FormControl'

const StyledSelect = withStyles(theme => ({
  root: {
    fontSize: theme.typography.fontSize,
  },
}))(props => <Select {...props} />)

export default function SelectSort({ options = [], variant, size, ...other }) {
  return (
    <FormControl {...{ variant, size }}>
      <StyledSelect id="select-sort" {...other}>
        {options.map(opt => (
          <MenuItem value={opt.value}>{opt.label}</MenuItem>
        ))}
      </StyledSelect>
    </FormControl>
  )
}
