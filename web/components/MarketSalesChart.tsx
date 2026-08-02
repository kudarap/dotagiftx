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

interface MarketSale {
  date: string
  avg: number
  count: number
}

interface ChartPoint {
  unix: number
  avg: number
  count: number
}

function formatDateUnix(unix: number) {
  return moment(unix).format('MMM D')
}

function formatXAxis(tickItem: number) {
  return formatDateUnix(tickItem)
}

interface CustomToolTipProps {
  active?: boolean
  payload?: Array<{ payload: ChartPoint }>
}

function CustomToolTip({ active, payload }: CustomToolTipProps) {
  if (!active) {
    return null
  }

  const p = payload?.[0]?.payload
  if (!p) {
    return null
  }

  return (
    <Paper style={{ padding: 8 }}>
      <strong>{formatDateUnix(p.unix)}</strong> <br />
      {amount(p.avg, 'USD')} <br />
      {p.count} sold
    </Paper>
  )
}

export default function MarketSalesChart({ data }: { data?: MarketSale[] }) {
  if (!data) {
    return null
  }

  const format: ChartPoint[] = data.map(v => ({
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
