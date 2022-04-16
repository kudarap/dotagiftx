import { makeStyles } from 'tss-react/mui'
import Link from '@/components/Link'

const badgeSettings = {
  supporter: {
    label: 'Supporter',
    color: '#596b95',
  },
  trader: {
    label: 'Trader',
    color: '#629cbd',
  },
  partner: {
    label: 'Partner',
    color: '#C79123',
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
