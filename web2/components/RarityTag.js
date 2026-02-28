import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@mui/material/Typography'
import {
  ITEM_RARITY_RARE,
  ITEM_RARITY_ULTRA_RARE,
  ITEM_RARITY_VERY_RARE,
} from '@/constants/palette'

// background: linear-gradient(#f9ffbf 10%, #fff 90%);
// text-shadow: 0px 0px 10px yellowgreen;
// WebkitBackgroundClip: text;
// WebkitTextFillColor: transparent;

// background: linear-gradient(#fdd08e 10%, #fff 90%);
// text-shadow: 0px 0px 10px darkorange;
// WebkitBackgroundClip: text;
// WebkitTextFillColor: transparent;

// background: linear-gradient(#F8E8B9 10%, #fff 90%);
// text-shadow: 0px 0px 10px goldenrod;
// WebkitBackgroundClip: text;
// WebkitTextFillColor: transparent;
const rarityStylerMap = {
  regular: null,
  rare: { color: ITEM_RARITY_RARE },
  'very rare': { color: ITEM_RARITY_VERY_RARE },
  'ultra rare': { color: ITEM_RARITY_ULTRA_RARE },
}

const getRarityStyle = value => {
  if (value === '') {
    return null
  }

  return { ...rarityStylerMap[value], textTransform: 'capitalize', display: 'inline' }
}

export default function RarityTag({ rarity, ...other }) {
  if (rarity === '' || rarity === 'regular') {
    return null
  }

  return (
    <Typography variant="caption" {...other} style={getRarityStyle(rarity)}>
      {` ${rarity}`}
    </Typography>
  )
}
RarityTag.propTypes = {
  rarity: PropTypes.string.isRequired,
}
