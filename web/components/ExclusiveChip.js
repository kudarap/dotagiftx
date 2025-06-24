import React from 'react'
import Link from '@mui/material/Link'
import Chip from '@mui/material/Chip'

export const tagSettings = {
  middleman: {
    label: 'Middleman',
    color: '#15803D',
    link: '/middleman',
  },
  moderator: {
    label: 'Moderator',
    color: '#9b59b6',
    link: '/moderators',
  },
}

export default function ExclusiveChip({ tag, ...props }) {
  if (!tag) {
    return null
  }

  const { label, color, link } = tagSettings[tag]
  return (
    <Chip
      size="small"
      variant="outlined"
      clickable
      component={Link}
      style={{
        textDecoration: 'none',
        backgroundColor: color,
        borderColor: color,
      }}
      label={label}
      href={link}
      {...props}
    />
  )
}
