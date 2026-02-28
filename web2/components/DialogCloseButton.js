import React from 'react'
import IconButton from '@mui/material/IconButton'
import CloseIcon from '@mui/icons-material/Close'

export default function DialogCloseButton(props) {
  return (
    <IconButton
      {...props}
      style={{ float: 'right' }}
      edge="start"
      color="inherit"
      aria-label="close"
      size="small">
      <CloseIcon />
    </IconButton>
  )
}
