import React from 'react'
import FormControlLabel from '@mui/material/FormControlLabel'
import Checkbox from '@mui/material/Checkbox'
import TextField from '@mui/material/TextField'
import { Typography } from '@mui/material'
import type { TextFieldProps } from '@mui/material/TextField'

export default function ReSellInput(props: TextFieldProps) {
  const [checked, setChecked] = React.useState(false)

  return (
    <div>
      <FormControlLabel
        style={{ color: '#ff9800' }}
        control={<Checkbox checked={checked} onChange={() => setChecked(!checked)} />}
        label={
          <Typography>
            <strong>Shopkeeper&apos;s Contract</strong>: I confirm this item exist on seller&apos;s
            inventory.
          </Typography>
        }
      />
      {checked && <TextField {...props} disabled={!checked} required={checked} autoFocus />}
      <br />
      <br />
    </div>
  )
}
