import React from 'react'
import Link from '@mui/material/Link'
import Chip from '@mui/material/Chip'

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
