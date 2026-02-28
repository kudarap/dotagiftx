import React from 'react'
import PropTypes from 'prop-types'
import MuiAvatar from '@mui/material/Avatar'
import { badgeSettings } from './SubscriberBadge'

const frameOptions = {
  donator: {
    border: 'goldenrod',
    frame: '/glow-frame.png',
  },
  aghanim: {
    border: '#4094ffed',
    frame: '/aghanim-frame.png',
  },
}

export default function Avatar(props) {
  const { glow, large = false, style: initStyle, src, badge, ...other } = props

  const glowFrame = frameOptions.donator
  let style = initStyle
  if (glow) {
    style = {
      ...style,
      border: `1px solid ${glowFrame.border}`,
      // animation: 'donatorglow4 12s infinite',
      // animationFillMode: 'forwards',
      // animationDelay: '3s',
      // animationTimingFunction: 'ease-in-out',
    }
  }
  if (badge) {
    const borderWidth = large ? 2 : 1
    style = {
      ...style,
      borderTop: `${borderWidth * 1}px solid ${badgeSettings[badge].color}`,
      borderLeft: `${borderWidth * 1}px solid ${badgeSettings[badge].color}`,
      borderRight: `${borderWidth * 1}px solid ${badgeSettings[badge].color}`,
      borderBottom: `${borderWidth * 2}px solid ${badgeSettings[badge].color}`,
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
          <img
            style={{ width: '100%', height: '100%', display: 'block' }}
            alt=""
            src={glowFrame.frame}
          />
        </div>
      )}
    </MuiAvatar>
  )
}

Avatar.propTypes = {
  style: PropTypes.object,
  glow: PropTypes.bool,
  large: PropTypes.bool,
  src: PropTypes.string,
  badge: PropTypes.string,
}

Avatar.defaultProps = {
  style: {},
  glow: false,
  large: false,
  src: null,
  badge: null,
}
