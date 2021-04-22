import React from 'react'
import PropTypes from 'prop-types'
import MuiAvatar from '@material-ui/core/Avatar'

export default function Avatar(props) {
  const { glow, style: initStyle, ...other } = props

  let style = initStyle
  if (glow) {
    style = {
      ...style,
      border: '1px solid goldenrod',
      animation: 'donatorglow4 12s infinite',
      animationFillMode: 'forwards',
      animationDelay: '3s',
      animationTimingFunction: 'ease-in-out',
    }
  }

  return <MuiAvatar style={style} {...other} />
}

Avatar.propTypes = {
  style: PropTypes.object,
  glow: PropTypes.bool,
}

Avatar.defaultProps = {
  style: {},
  glow: false,
}
