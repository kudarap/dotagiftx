import React, { useEffect, useRef } from 'react'
import PropTypes from 'prop-types'

// Hero background video, rendered on top of the static hero image. The video
// only starts downloading after the user's first interaction (gesture), so it
// never competes with the LCP image or delays first paint.
function HeroVideo({ sources, className, style }) {
  const videoRef = useRef(null)

  useEffect(() => {
    const video = videoRef.current
    if (!video) {
      return
    }

    let started = false
    const startPlayback = () => {
      if (started) {
        return
      }
      started = true
      video.play().catch(() => {
        started = false
      })
    }

    const onGesture = () => startPlayback()
    const options = { passive: true }
    window.addEventListener('pointerdown', onGesture, options)
    window.addEventListener('keydown', onGesture, options)

    return () => {
      window.removeEventListener('pointerdown', onGesture)
      window.removeEventListener('keydown', onGesture)
    }
  }, [])

  const mimeType = url => {
    const ext = url.split('.').pop().split('?')[0]
    if (ext === 'mp4') {
      return 'video/mp4'
    }
    if (ext === 'webm') {
      return 'video/webm'
    }
    return undefined
  }

  return (
    <video
      ref={videoRef}
      className={className}
      style={style}
      preload="none"
      muted
      loop
      playsInline>
      {sources.map(src => (
        <source key={src} type={mimeType(src)} src={src} />
      ))}
    </video>
  )
}

HeroVideo.propTypes = {
  sources: PropTypes.arrayOf(PropTypes.string).isRequired,
  className: PropTypes.string,
  style: PropTypes.object,
}

HeroVideo.defaultProps = {
  className: '',
  style: {},
}

export default HeroVideo
