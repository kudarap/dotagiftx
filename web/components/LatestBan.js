import moment from 'moment'
import React from 'react'

import { blacklistSearch } from '@/service/api'

export default function LatestBan() {
  const [ban, setBan] = React.useState(null)

  React.useEffect(() => {
    ;(async () => {
      try {
        const user = await blacklistSearch({ limit: 1 })
        setBan(user[0] || null)
      } catch (error) {
        console.log('error getting lastest ban', error)
      }
    })()
  }, [])

  if (!ban) {
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
      {moment(ban.updated_at).fromNow()}
    </span>
  )
}
