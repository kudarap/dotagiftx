import React from 'react'
import Link from '@material-ui/core/Link'
import Chip from '@material-ui/core/Chip'

export default function ChipLink(props) {
  return (
    <Chip
      size="small"
      variant="outlined"
      clickable
      component={Link}
      target="_blank"
      rel="noreferrer noopener"
      style={{ textDecoration: 'none' }}
      {...props}
    />
  )
}
