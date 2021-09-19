import moment from 'moment'
import React from 'react'

import AppContext from './AppContext'

export default function LatestBan() {
  const { latestBan } = React.useContext(AppContext)
  if (!latestBan) {
    return null
  }

  return (
    <span
      style={{
        position: 'absolute',
        fontSize: '0.6rem',
        display: 'block',
        marginTop: '-3px',
        color: '#FF6464',
      }}>
      {moment(latestBan.updated_at).fromNow()}
    </span>
  )
}
