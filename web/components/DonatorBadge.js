import { makeStyles } from '@material-ui/core/styles'
import Link from '@/components/Link'

const useStyles = makeStyles(theme => ({
  root: {
    color: 'white',
    padding: '2px 4px 0px',
    fontSize: 10,
    background: 'goldenrod',
    fontWeight: 500,
    borderRadius: '2px',
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