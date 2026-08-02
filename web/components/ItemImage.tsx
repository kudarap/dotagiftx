import React from 'react'
import Image from 'next/image'
import type { CSSProperties } from 'react'
import Box from '@mui/material/Box'
import { CDN_URL } from '@/service/api'
import { itemRarityColorMap } from '@/constants/palette'

const baseSizeQuality = 20
export function retinaSrcSet(filename: string, width: number, height: number) {
  if (!filename) {
    return { src: '' }
  }

  const src = `${CDN_URL}/${width + baseSizeQuality}x${height + baseSizeQuality}/${filename}`
  const src2x = `${CDN_URL}/${width * 2}x${height * 2}/${filename}`
  return { src, srcSet: `${src} 1x, ${src2x} 2x` }
}

interface ItemImageProps {
  image: string
  width?: number
  height?: number
  title?: string
  rarity?: string
  className?: string
  style?: CSSProperties
}

export default function ItemImage({ image, title, rarity, className, width, height, ...other }: ItemImageProps) {
  const containerStyle: CSSProperties = {
    lineHeight: 1,
    flexShrink: 0,
    overflow: 'hidden',
    userSelect: 'none',
    display: 'flex',
    flexDirection: 'column',
  }
  if (rarity) {
    containerStyle.border = `1px solid ${itemRarityColorMap[rarity]}`
  }

  const imageStyle: CSSProperties = {
    color: 'transparent',
    width: '100%',
    height: 'auto',
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
        priority
      />
    </Box>
  )
}
