import React from 'react'
import PropTypes from 'prop-types'
import Image from 'next/image'
import Box from '@mui/material/Box'
import { CDN_URL } from '@/service/api'
import { itemRarityColorMap, itemRarityGlowSet } from '@/constants/palette'

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
  preload,
  ...other
}) {
  const containerStyle = {
    lineHeight: 1,
    flexShrink: 0,
    overflow: 'hidden',
    userSelect: 'none',
    display: 'flex',
    flexDirection: 'column',
  }
  if (rarity) {
    const color = itemRarityColorMap[rarity]
    containerStyle.border = `1px solid ${color}`
    if (itemRarityGlowSet.includes(rarity)) {
      containerStyle.boxShadow = `0 0 8px ${color}`
    }
  }

  const imageStyle = {
    color: 'transparent',
    width: '100%',
    height: '100%',
    objectFit: 'cover',
  }

  let baseSrc = CDN_URL + image
  // using srcset to support high dpi or retina displays when
  // dimension were set.
  if (width && height) {
    const rs = retinaSrcSet(image, width, height)
    baseSrc = rs.src
  }

  return (
    <Box style={containerStyle} className={className} {...other}>
      <Image
        style={imageStyle}
        src={baseSrc}
        alt={title || image}
        width={width}
        height={height}
        quality={100}
        preload={preload}
      />
    </Box>
  )
}
ItemImage.propTypes = {
  image: PropTypes.string.isRequired,
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  title: PropTypes.string,
  rarity: PropTypes.string,
  className: PropTypes.string,
  preload: PropTypes.bool,
}
ItemImage.defaultProps = {
  title: null,
  rarity: null,
  className: '',
  preload: false,
}
