import React from 'react'
import { styled } from '@mui/material/styles'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import type { SelectProps } from '@mui/material/Select'
import FormControl from '@mui/material/FormControl'

interface SelectOption {
  value: string | number
  label: string
}

type SelectSortProps = SelectProps & {
  options: SelectOption[]
}

const StyledSelect = styled(Select)(({ theme }) => ({
  fontSize: theme.typography.fontSize,
}))

export default function SelectSort({ options, variant, size, ...other }: SelectSortProps) {
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
