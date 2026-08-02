import React from 'react'
import MuiLink from '@mui/material/Link'
import Chip from '@mui/material/Chip'
import type { ChipProps } from '@mui/material/Chip'

export default function ChipLink(props: ChipProps & { href?: string }) {
  return (
    <Chip
      size="small"
      variant="outlined"
      clickable
      component={MuiLink}
      target="_blank"
      rel="noreferrer noopener"
      style={{ textDecoration: 'none' }}
      {...props}
    />
  )
}
