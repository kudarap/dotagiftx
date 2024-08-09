import React from 'react'
import PropTypes from 'prop-types'
import { CDN_URL } from '@/service/api'
import { itemRarityColorMap } from '@/constants/palette'
import Image from 'next/image'

const baseSizeQuality = 20
export function retinaSrcSet(filename, width, height) {
  if (!filename) {
    return { src: '' }
  }

  const src = `${CDN_URL}/${width + baseSizeQuality}x${height + baseSizeQuality}/${filename}`
  const src2x = `${CDN_URL}/${width * 2}x${height * 2}/${filename}`
  return { src, srcSet: `${src} 1x, ${src2x} 2x` }
}

export default function ItemImage({
  image,
  title,
  rarity,
  className,
  width,
  height,
  nextOptimized,
  ...other
}) {
  const contStyle = {
    lineHeight: 1,
    flexShrink: 0,
    overflow: 'hidden',
    userSelect: 'none',
  }

  if (rarity) {
    contStyle.border = `1px solid ${itemRarityColorMap[rarity]}`
  }

  const imgStyle = {
    color: 'transparent',
    objectFit: 'cover',
    textAlign: 'center',
    textIndent: '10000px',
  }

  let baseSrc = CDN_URL + image
  // using srcset to support high dpi or retina displays when
  // dimension were set.
  let srcSet = null
  if (width && height) {
    const rs = retinaSrcSet(image, width, height)
    baseSrc = rs.src
    srcSet = rs.srcSet
  }

  if (!nextOptimized) {
    return (
      <div style={contStyle} className={className}>
        <img
          loading="lazy"
          src={baseSrc}
          alt={title || image}
          style={imgStyle}
          height={height}
          {...other}
        />
      </div>
    )
  }

  return (
    <div style={contStyle} className={className}>
      <Image
        src={baseSrc}
        alt={title || image}
        style={imgStyle}
        width={width}
        height={height}
        quality={100}
        priority
        responsive
        {...other}
      />
    </div>
  )
}
ItemImage.propTypes = {
  image: PropTypes.string.isRequired,
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  title: PropTypes.string,
  rarity: PropTypes.string,
  className: PropTypes.string,
  nextOptimized: PropTypes.bool,
}
ItemImage.defaultProps = {
  title: null,
  rarity: null,
  className: '',
  nextOptimized: false,
}
