import { makeStyles } from 'tss-react/mui'
import type { CSSProperties } from 'react'
import Link from '@/components/Link'

const useStyles = makeStyles()(() => ({
  root: {
    color: 'white',
    padding: '0 0.675rem',
    fontSize: 10,
    background: 'goldenrod',
    fontWeight: '0.875rem',
    borderRadius: '2px',
    border: '1px solid goldenrod',
    display: 'inline-block',
  },
}))

interface DonatorBadgeProps {
  style?: CSSProperties
  size?: 'medium' | 'large'
  className?: string
}

export default function DonatorBadge({ style: initialStyle, size, ...other }: DonatorBadgeProps) {
  const { classes } = useStyles()

  const currentStyle: CSSProperties = { ...initialStyle }
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
