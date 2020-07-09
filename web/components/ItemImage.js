import React from 'react'
import PropTypes from 'prop-types'
import { CDN_URL } from '@/service/api'
import { itemRarityColorMap } from '../constants/palette'

export default function ItemImage({ image, title, rarity, ...other }) {
  const contStyle = {
    display: 'flex',
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
    width: '100%',
    height: '100%',
    objectFit: 'cover',
    textAlign: 'center',
    textIndent: '10000px',
  }

  return (
    <div style={contStyle} {...other}>
      <img src={CDN_URL + image} alt={title || image} style={imgStyle} />
    </div>
  )
}
ItemImage.propTypes = {
  image: PropTypes.string.isRequired,
  title: PropTypes.string,
  rarity: PropTypes.string,
}
ItemImage.defaultProps = {
  title: null,
  rarity: null,
}
