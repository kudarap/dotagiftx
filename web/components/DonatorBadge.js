import { makeStyles } from '@material-ui/core/styles'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  root: {
    color: 'white',
    padding: '0 0.675rem',
    fontSize: 10,
    background: 'linear-gradient(to right, #1A201F, #312d26)',
    fontWeight: '0.875rem',
    borderRadius: '2px',
    border: '1px solid goldenrod',
    display: 'inline-block',
  },
}))

export default function DonatorBadge({ style: initialStyle, size, ...other }) {
  const classes = useStyles()

  const currentStyle = { ...initialStyle }
  if (size === 'medium') {
    currentStyle.fontSize = '0.875rem'
  }

  return (
    <Link
      className={classes.root}
      style={currentStyle}
      {...other}
      href="/donate"
      underline="none"
    />
  )
}

DonatorBadge.defaultProps = {
  style: {},
  size: false,
}
