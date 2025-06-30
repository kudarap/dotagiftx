import React from 'react'
import Link from '@mui/material/Link'
import Chip from '@mui/material/Chip'
import {
  USER_SUBSCRIPTION_MAP_COLOR,
  USER_SUBSCRIPTION_PARTNER,
  USER_SUBSCRIPTION_SUPPORTER,
  USER_SUBSCRIPTION_TRADER,
} from '@/constants/user'

export const tagSettings = {
  // subscribers
  supporter: {
    label: 'Supporter',
    color: USER_SUBSCRIPTION_MAP_COLOR[USER_SUBSCRIPTION_SUPPORTER],
    link: '/plus',
  },
  trader: {
    label: 'Trader',
    color: USER_SUBSCRIPTION_MAP_COLOR[USER_SUBSCRIPTION_TRADER],
    link: '/plus',
  },
  partner: {
    label: 'Partner',
    color: USER_SUBSCRIPTION_MAP_COLOR[USER_SUBSCRIPTION_PARTNER],
    link: '/plus',
  },
  // internals
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
