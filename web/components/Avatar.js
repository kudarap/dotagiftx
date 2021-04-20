import React from 'react'
import PropTypes from 'prop-types'
import clsx from 'clsx'
import { makeStyles } from '@material-ui/core/styles'
import MuiAvatar from '@material-ui/core/Avatar'

const useStyles = makeStyles({
  glow: {
    border: '1px solid goldenrod',
    animation: 'donatorglow4 13s infinite',
    animationFillMode: 'forwards',
    animationDelay: '3s',
    animationTimingFunction: 'ease-in-out',
  },
})

export default function Avatar(props) {
  const classes = useStyles()

  const { glow, className: classNameProps, ...other } = props
  const className = clsx(classNameProps, glow ? classes.glow : null)

  return <MuiAvatar className={className} {...other} />
}

Avatar.propTypes = {
  className: PropTypes.object,
  glow: PropTypes.bool,
}

Avatar.defaultProps = {
  className: {},
  glow: false,
}
