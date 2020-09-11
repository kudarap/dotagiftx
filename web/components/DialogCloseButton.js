import React from 'react'
import IconButton from '@material-ui/core/IconButton'
import CloseIcon from '@material-ui/icons/Close'

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
