import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import {
  ITEM_RARITY_RARE,
  ITEM_RARITY_ULTRA_RARE,
  ITEM_RARITY_VERY_RARE,
} from '@/constants/palette'

// background: linear-gradient(#f9ffbf 10%, #fff 90%);
// text-shadow: 0px 0px 10px yellowgreen;
// -webkit-background-clip: text;
// -webkit-text-fill-color: transparent;

// background: linear-gradient(#fdd08e 10%, #fff 90%);
// text-shadow: 0px 0px 10px darkorange;
// -webkit-background-clip: text;
// -webkit-text-fill-color: transparent;

// background: linear-gradient(#F8E8B9 10%, #fff 90%);
// text-shadow: 0px 0px 10px goldenrod;
// -webkit-background-clip: text;
// -webkit-text-fill-color: transparent;
const rarityStylerMap = {
  regular: null,
  rare: { color: ITEM_RARITY_RARE },
  'very rare': { color: ITEM_RARITY_VERY_RARE },
  'ultra rare': { color: ITEM_RARITY_ULTRA_RARE },
}
const rarityStylerMap2 = {
  regular: null,
  rare: {
    color: 'yellowgreen',
    background: 'linear-gradient(#f9ffbf 10%, #fff 90%)',
    '-webkit-background-clip': 'text',
    '-webkit-text-fill-color': 'transparent',
    filter: 'drop-shadow(0px 0px 2px yellowgreen)',
  },
  'very rare': {
    background: 'linear-gradient(#fdd08e 10%, #fff 90%)',
    '-webkit-background-clip': 'text',
    '-webkit-text-fill-color': 'transparent',
    filter: 'drop-shadow(0px 0px 2px darkorange)',
  },
  'ultra rare': {
    background: 'linear-gradient(#F8E8B9 10%, #fff 90%)',
    '-webkit-background-clip': 'text',
    '-webkit-text-fill-color': 'transparent',
    filter: 'drop-shadow(0px 0px 2px goldenrod)',
  },
}
const getRarityStyle = value => {
  if (value === '') {
    return null
  }

  return rarityStylerMap[value]
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
