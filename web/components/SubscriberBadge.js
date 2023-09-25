import { makeStyles } from 'tss-react/mui'
import Link from '@/components/Link'
import {
  USER_SUBSCRIPTION_MAP_COLOR,
  USER_SUBSCRIPTION_PARTNER,
  USER_SUBSCRIPTION_SUPPORTER,
  USER_SUBSCRIPTION_TRADER,
} from '@/constants/user'

export const badgeSettings = {
  middleman: {
    label: 'Middleman',
    color: '#15803D',
  },
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

const useStyles = makeStyles()(theme => ({
  root: {
    color: 'white',
    padding: '0 4px',
    padding: '0 0.675rem',
    fontSize: 10,
    fontWeight: 500,
    fontWeight: '0.875rem',
    borderRadius: '2px',
    display: 'inline-block',
    textTransform: 'uppercase',
    border: '1px solid gray',
  },
}))

export default function SubscriberBadge({ style: initialStyle, size, type, ...other }) {
  const { classes } = useStyles()

  const currentStyle = { ...initialStyle }
  if (size === 'medium') {
    currentStyle.fontSize = '0.875rem'
  }
  if (type) {
    currentStyle.background = badgeSettings[type].color
    currentStyle.borderColor = badgeSettings[type].color
  }

  return (
    <Link className={classes.root} style={currentStyle} {...other} href="/plus" underline="none">
      {badgeSettings[type].label}
    </Link>
  )
}
SubscriberBadge.defaultProps = {
  style: {},
  size: false,
  type: '',
}
