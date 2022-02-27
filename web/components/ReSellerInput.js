import React from 'react'
import FormControlLabel from '@mui/material/FormControlLabel'
import Checkbox from '@mui/material/Checkbox'
import TextField from '@mui/material/TextField'

export default function ReSellInput(props) {
  const [checked, setChecked] = React.useState(false)

  return (
    <div>
      <FormControlLabel
        style={{ color: '#ff9800' }}
        control={<Checkbox checked={checked} onChange={() => setChecked(!checked)} />}
        label="Item Resell - I confirm this item exist on seller's inventory."
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
