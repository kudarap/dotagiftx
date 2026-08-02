import { makeStyles } from 'tss-react/mui'
import type { CSSProperties } from 'react'
import Link from '@/components/Link'
import {
  USER_SUBSCRIPTION_MAP_COLOR,
  USER_SUBSCRIPTION_PARTNER,
  USER_SUBSCRIPTION_SUPPORTER,
  USER_SUBSCRIPTION_TRADER,
} from '@/constants/user'

export const badgeSettings = {
  supporter: {
    label: 'Supporter',
    color: USER_SUBSCRIPTION_MAP_COLOR[USER_SUBSCRIPTION_SUPPORTER],
  },
  trader: {
    label: 'Trader',
    color: USER_SUBSCRIPTION_MAP_COLOR[USER_SUBSCRIPTION_TRADER],
  },
  partner: {
    label: 'Partner',
    color: USER_SUBSCRIPTION_MAP_COLOR[USER_SUBSCRIPTION_PARTNER],
  },
}

export type SubscriberBadgeType = keyof typeof badgeSettings

const useStyles = makeStyles()(() => ({
  root: {
    color: 'white',
    padding: '0 0.675rem',
    fontSize: 10,
    fontWeight: '0.875rem',
    borderRadius: '2px',
    display: 'inline-block',
    textTransform: 'uppercase',
    border: '1px solid gray',
  },
}))

interface SubscriberBadgeProps {
  style?: CSSProperties
  size?: 'medium' | 'large'
  type?: SubscriberBadgeType
  className?: string
}

export default function SubscriberBadge({ style: initialStyle, size, type, ...other }: SubscriberBadgeProps) {
  const { classes } = useStyles()

  const currentStyle: CSSProperties = { ...initialStyle }
  if (size === 'medium') {
    currentStyle.fontSize = '0.875rem'
  }
  if (size === 'large') {
    currentStyle.fontSize = '1rem'
  }
  if (type) {
    currentStyle.background = badgeSettings[type].color
    currentStyle.borderColor = badgeSettings[type].color
  }

  return (
    <Link className={classes.root} style={currentStyle} {...other} href="/plus" underline="none">
      {type ? badgeSettings[type].label : null}
    </Link>
  )
}
