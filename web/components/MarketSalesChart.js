import React from 'react'
import PropTypes from 'prop-types'
import { format, getUnixTime } from 'date-fns'
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import Paper from '@mui/material/Paper'
import { lightGreen as graphColor } from '@mui/material/colors'

import { amount } from '@/lib/format'

function formatDateUnix(unix) {
  return format(new Date(unix), 'MMM d')
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

CustomToolTip.propTypes = {
  active: PropTypes.bool,
  payload: PropTypes.array.isRequired,
}

CustomToolTip.defaultProps = {
  active: false,
}

export default function MarketSalesChart({ data }) {
  if (!data) {
    return null
  }

  const format = data.map(v => ({
    unix: getUnixTime(new Date(v.date) * 1000),
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

MarketSalesChart.propTypes = {
  data: PropTypes.array,
}

MarketSalesChart.defaultProps = {
  data: undefined,
}
