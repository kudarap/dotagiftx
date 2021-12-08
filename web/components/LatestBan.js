import React, { useContext, useMemo, useState } from 'react'
import moment from 'moment'

import AppContext from './AppContext'

const sinceDayMin = 1
const sinceDayMax = 30
const sinceRate = sinceDayMax / sinceDayMin

const getDaysFromTs = datetime => {
  const ts = moment().diff(datetime)
  return Math.ceil(ts / 86400000)
}

export default function LatestBan() {
  const { latestBan } = useContext(AppContext)
  const [grayscale, setGrayscale] = useState(0)

  const recentBanAt = latestBan?.updated_at || null

  useMemo(() => {
    if (!recentBanAt) {
      return
    }

    const daysDiff = getDaysFromTs(recentBanAt)
    setGrayscale((daysDiff / sinceRate).toFixed(2) * 100)
  }, [recentBanAt])

  if (!recentBanAt) {
    return null
  }

  return (
    <span
      style={{
        position: 'absolute',
        fontSize: '0.6rem',
        display: 'block',
        marginTop: '-5px',
        color: '#FF6464',
        filter: `grayscale(${grayscale}%)`,
      }}>
      {moment(recentBanAt).fromNow()}
    </span>
  )
}
