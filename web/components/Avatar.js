import React from 'react'
import PropTypes from 'prop-types'
import MuiAvatar from '@material-ui/core/Avatar'

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

  // return (
  //   <div style={style}>
  //     <img alt="" {...other} />
  //     <div style={{ position: 'absolute', margin: -12 }}>
  //       <img style={{ width: '100%', height: '100%', display: 'block' }} alt="" src="/glow.png" />
  //     </div>
  //   </div>
  // )
  // return <MuiAvatar style={style} {...other} />
}

Avatar.propTypes = {
  style: PropTypes.object,
  glow: PropTypes.bool,
}

Avatar.defaultProps = {
  style: {},
  glow: false,
}
