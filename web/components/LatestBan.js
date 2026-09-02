import React, { useState } from 'react'
import { differenceInMilliseconds, parseISO } from 'date-fns'
import { relativeFromNow } from '@/lib/format'
import { blacklistSearch } from '@/service/api'
import { get, save } from '@/service/storage'

const sinceDayMin = 1
const sinceDayMax = 30
const sinceRate = sinceDayMax / sinceDayMin

const getDaysFromTs = datetime => {
  const ts = differenceInMilliseconds(new Date(), parseISO(datetime))
  return Math.ceil(ts / 86400000)
}

const CACHE_KEY = 'latest-ban'

export default function LatestBan() {
  const [value, setValue] = useState(get(CACHE_KEY))

  React.useEffect(() => {
    ;(async () => {
      // skip get
      if (value) {
        return
      }

      try {
        const res = await blacklistSearch({ limit: 1, sort: 'updated_at:desc' })
        if (res && res.length > 0) {
          const latest = res[0].updated_at
          setValue(latest)
          save(CACHE_KEY, latest, 3600)
        }
      } catch (error) {
        console.warn('failed getting lastest ban', error)
      }
    })()
  }, [value])

  if (!value) {
    return null
  }

  const grayscale = (getDaysFromTs(value) / sinceRate).toFixed(2) * 100 || 0
  return (
    <span
      style={{
        position: 'absolute',
        fontSize: '0.55rem',
        display: 'block',
        marginTop: '-0.16rem',
        color: '#FF6464',
        filter: `grayscale(${grayscale}%)`,
      }}
    >
      {relativeFromNow(value)}
    </span>
  )
}
