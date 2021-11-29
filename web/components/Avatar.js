import React from 'react'
import PropTypes from 'prop-types'
import MuiAvatar from '@mui/material/Avatar'

export default function Avatar(props) {
  const { glow, style: initStyle, src, ...other } = props

  let style = initStyle
  if (glow) {
    style = {
      ...style,
      border: '1px solid goldenrod',
      // animation: 'donatorglow4 12s infinite',
      // animationFillMode: 'forwards',
      // animationDelay: '3s',
      // animationTimingFunction: 'ease-in-out',
    }
  }

  if (!glow) {
    return <MuiAvatar src={src} style={style} {...other} />
  }

  return (
    <MuiAvatar style={style} {...other}>
      <img src={src} alt="" style={{ width: '100%', height: '100%' }} />
      {glow && (
        <div style={{ position: 'absolute', margin: '-12%' }}>
          <img style={{ width: '100%', height: '100%', display: 'block' }} alt="" src="/glow.png" />
        </div>
      )}
    </MuiAvatar>
  )
}

Avatar.propTypes = {
  style: PropTypes.object,
  glow: PropTypes.bool,
  src: PropTypes.string,
}

Avatar.defaultProps = {
  style: {},
  glow: false,
  src: null,
}
