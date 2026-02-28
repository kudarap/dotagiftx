import React from 'react'
import moment from 'moment'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from 'recharts'
import Paper from '@mui/material/Paper'
import { lightGreen as graphColor } from '@mui/material/colors'

import { amount } from '@/lib/format'

function formatDateUnix(unix) {
  return moment(unix).format('MMM D')
}

function formatXAxis(tickItem) {
  return formatDateUnix(tickItem)
}

function CustomToolTip(props) {
  const { active } = props
  if (!active) {
    return null
  }

  const { payload } = props
  const p = payload[0].payload
  return (
    <Paper style={{ padding: 8 }}>
      <strong>{formatDateUnix(p.unix)}</strong> <br />
      {amount(p.avg, 'USD')} <br />
      {p.count} sold
    </Paper>
  )
}

export default function MarketSalesChart({ data }) {
  if (!data) {
    return null
  }

  const format = data.map(v => ({
    unix: moment(v.date).unix() * 1000,
    avg: Number(v.avg.toFixed(2)),
    count: v.count,
  }))

  return (
    <div style={{ width: '100%', height: 200, marginLeft: -20 }}>
      <ResponsiveContainer>
        <LineChart data={format}>
          <CartesianGrid strokeDasharray="3 3" stroke="#555" />
          <XAxis
            dataKey="unix"
            type="number"
            domain={['dataMin', 'dataMax']}
            tickFormatter={formatXAxis}
          />
          <YAxis />
          <Legend />
          <Tooltip content={<CustomToolTip />} />
          <Line
            name="Average Sale Prices"
            type="linear"
            dataKey="avg"
            stroke={graphColor[800]}
            dot={false}
            strokeWidth={2}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
