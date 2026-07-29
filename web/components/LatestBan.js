import React, { useState } from 'react'
import moment from 'moment'

const sinceDayMin = 1
const sinceDayMax = 30
const sinceRate = sinceDayMax / sinceDayMin

const getDaysFromTs = datetime => {
  const ts = moment().diff(datetime)
  return Math.ceil(ts / 86400000)
}

export default function LatestBan({ value }) {
  const [grayscale, setGrayscale] = useState(0)

  if (!value) {
    return null
  }

  const recentBanAt = getDaysFromTs(value)
  setGrayscale((recentBanAt / sinceRate).toFixed(2) * 100)

  return (
    <span
      style={{
        position: 'absolute',
        fontSize: '0.6rem',
        display: 'block',
        marginTop: '-0.24rem',
        color: '#FF6464',
        filter: `grayscale(${grayscale}%)`,
      }}>
      {moment(recentBanAt).fromNow()}
    </span>
  )
}
