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
import graphColor from '@mui/material/colors/lightGreen'
import { amount } from '@/lib/format'

const testdata = JSON.parse(
  '[{"date":"2020-10-15T00:00:00Z","avg":7,"count":1},{"date":"2020-10-28T00:00:00Z","avg":10,"count":1},{"date":"2020-11-05T00:00:00Z","avg":10,"count":1},{"date":"2020-11-19T00:00:00Z","avg":12.5,"count":1}]'
  // '[{"avg":13.7,"count":5,"date":{"$reql_type$":"TIME","epoch_time":1600041600,"timezone":"+00:00"}},{"avg":82,"count":1,"date":{"$reql_type$":"TIME","epoch_time":1602633600,"timezone":"+00:00"}},{"avg":7.1,"count":10,"date":{"$reql_type$":"TIME","epoch_time":1602720000,"timezone":"+00:00"}},{"avg":15,"count":1,"date":{"$reql_type$":"TIME","epoch_time":1603584000,"timezone":"+00:00"}},{"avg":10,"count":3,"date":{"$reql_type$":"TIME","epoch_time":1603843200,"timezone":"+00:00"}},{"avg":1.5,"count":1,"date":{"$reql_type$":"TIME","epoch_time":1604361600,"timezone":"+00:00"}},{"avg":10,"count":1,"date":{"$reql_type$":"TIME","epoch_time":1604534400,"timezone":"+00:00"}},{"avg":9,"count":4,"date":{"$reql_type$":"TIME","epoch_time":1604793600,"timezone":"+00:00"}},{"avg":4.25,"count":4,"date":{"$reql_type$":"TIME","epoch_time":1604880000,"timezone":"+00:00"}},{"avg":4.5,"count":2,"date":{"$reql_type$":"TIME","epoch_time":1605225600,"timezone":"+00:00"}},{"avg":10.75,"count":2,"date":{"$reql_type$":"TIME","epoch_time":1605312000,"timezone":"+00:00"}},{"avg":8.666666666666666,"count":3,"date":{"$reql_type$":"TIME","epoch_time":1605657600,"timezone":"+00:00"}},{"avg":25.9,"count":5,"date":{"$reql_type$":"TIME","epoch_time":1605744000,"timezone":"+00:00"}},{"avg":3.25,"count":2,"date":{"$reql_type$":"TIME","epoch_time":1605916800,"timezone":"+00:00"}},{"avg":28.7,"count":5,"date":{"$reql_type$":"TIME","epoch_time":1606003200,"timezone":"+00:00"}},{"avg":1.5,"count":4,"date":{"$reql_type$":"TIME","epoch_time":1606176000,"timezone":"+00:00"}},{"avg":3.375,"count":6,"date":{"$reql_type$":"TIME","epoch_time":1606262400,"timezone":"+00:00"}},{"avg":21.25,"count":2,"date":{"$reql_type$":"TIME","epoch_time":1606348800,"timezone":"+00:00"}},{"avg":2.8333333333333335,"count":3,"date":{"$reql_type$":"TIME","epoch_time":1606435200,"timezone":"+00:00"}},{"avg":8.833333333333334,"count":3,"date":{"$reql_type$":"TIME","epoch_time":1606521600,"timezone":"+00:00"}},{"avg":5.1,"count":5,"date":{"$reql_type$":"TIME","epoch_time":1606608000,"timezone":"+00:00"}},{"avg":4.25,"count":2,"date":{"$reql_type$":"TIME","epoch_time":1606780800,"timezone":"+00:00"}},{"avg":14.25,"count":2,"date":{"$reql_type$":"TIME","epoch_time":1606867200,"timezone":"+00:00"}},{"avg":1.6666666666666667,"count":3,"date":{"$reql_type$":"TIME","epoch_time":1606953600,"timezone":"+00:00"}},{"avg":32,"count":1,"date":{"$reql_type$":"TIME","epoch_time":1607040000,"timezone":"+00:00"}},{"avg":2.125,"count":4,"date":{"$reql_type$":"TIME","epoch_time":1607126400,"timezone":"+00:00"}},{"avg":1.5,"count":1,"date":{"$reql_type$":"TIME","epoch_time":1607299200,"timezone":"+00:00"}}]'
)

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

  const format = data.map(v => {
    return {
      unix: moment(v.date).unix() * 1000,
      avg: Number(v.avg.toFixed(2)),
      count: v.count,
    }
  })

  return (
    <div style={{ width: '100%', height: 200 }}>
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
