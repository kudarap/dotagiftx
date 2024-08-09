import React from 'react'
import FormControlLabel from '@mui/material/FormControlLabel'
import Checkbox from '@mui/material/Checkbox'
import TextField from '@mui/material/TextField'
import { Typography } from '@mui/material'

export default function ReSellInput(props) {
  const [checked, setChecked] = React.useState(false)

  return (
    <div>
      <FormControlLabel
        style={{ color: '#ff9800' }}
        control={<Checkbox checked={checked} onChange={() => setChecked(!checked)} />}
        label={
          <Typography>
            <strong>Shopkeeper's Contract</strong>: I confirm this item exist on seller's inventory.
          </Typography>
        }
      />
      {checked && (
        <>
          <TextField {...props} disabled={!checked} required={checked} autoFocus />
        </>
      )}
      <br />
      <br />
    </div>
  )
}
