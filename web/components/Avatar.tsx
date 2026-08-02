import React from 'react'
import Image from 'next/image'
import type { CSSProperties } from 'react'
import MuiAvatar from '@mui/material/Avatar'
import type { AvatarProps as MuiAvatarProps } from '@mui/material/Avatar'
import { badgeSettings } from './SubscriberBadge'
import type { SubscriberBadgeType } from './SubscriberBadge'

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

interface AvatarProps extends MuiAvatarProps {
  glow?: boolean
  large?: boolean
  src?: string
  badge?: SubscriberBadgeType
  href?: string
}

export default function Avatar(props: AvatarProps) {
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
    return <MuiAvatar src={src} style={style as CSSProperties} {...other} />
  }

  return (
    <MuiAvatar style={style as CSSProperties} {...other}>
      <Image src={src as string} alt="" style={{ width: '100%', height: '100%' }} />
      {glow && (
        <div style={{ position: 'absolute', margin: '-12%' }}>
          <Image
            style={{ width: '100%', height: '100%', display: 'block' }}
            alt=""
            src={glowFrame.frame}
          />
        </div>
      )}
    </MuiAvatar>
  )
}
