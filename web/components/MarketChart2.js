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
import Paper from '@material-ui/core/Paper'
import primary from '@material-ui/core/colors/lightGreen'
import { amount } from '@/lib/format'

const testdata = JSON.parse(
  '[{"avg":13.7,"count":5,"day":14,"month":9,"year":2020},{"avg":82,"count":1,"day":14,"month":10,"year":2020},{"avg":7.1,"count":10,"day":15,"month":10,"year":2020},{"avg":15,"count":1,"day":25,"month":10,"year":2020},{"avg":10,"count":3,"day":28,"month":10,"year":2020},{"avg":1.5,"count":1,"day":3,"month":11,"year":2020},{"avg":10,"count":1,"day":5,"month":11,"year":2020},{"avg":9,"count":4,"day":8,"month":11,"year":2020},{"avg":4.25,"count":4,"day":9,"month":11,"year":2020},{"avg":4.5,"count":2,"day":13,"month":11,"year":2020},{"avg":10.75,"count":2,"day":14,"month":11,"year":2020},{"avg":8.666666666666666,"count":3,"day":18,"month":11,"year":2020},{"avg":25.9,"count":5,"day":19,"month":11,"year":2020},{"avg":3.25,"count":2,"day":21,"month":11,"year":2020},{"avg":28.7,"count":5,"day":22,"month":11,"year":2020},{"avg":1.5,"count":4,"day":24,"month":11,"year":2020},{"avg":3.375,"count":6,"day":25,"month":11,"year":2020},{"avg":21.25,"count":2,"day":26,"month":11,"year":2020},{"avg":2.8333333333333335,"count":3,"day":27,"month":11,"year":2020},{"avg":8.833333333333334,"count":3,"day":28,"month":11,"year":2020},{"avg":5.1,"count":5,"day":29,"month":11,"year":2020},{"avg":4.25,"count":2,"day":1,"month":12,"year":2020},{"avg":14.25,"count":2,"day":2,"month":12,"year":2020},{"avg":1.6666666666666667,"count":3,"day":3,"month":12,"year":2020},{"avg":32,"count":1,"day":4,"month":12,"year":2020},{"avg":2.125,"count":4,"day":5,"month":12,"year":2020},{"avg":1.5,"count":1,"day":7,"month":12,"year":2020}]'
  // '[{"avg":1.5,"count":1,"day":8,"month":11,"year":2020},{"avg":1.5,"count":2,"day":24,"month":11,"year":2020},{"avg":2.5,"count":1,"day":26,"month":11,"year":2020},{"avg":1.5,"count":1,"day":5,"month":12,"year":2020},{"avg":1.5,"count":1,"day":7,"month":12,"year":2020}]'
).map(v => {
  const d = new Date(v.year, v.month, v.day)
  return {
    unix: d.getTime(),
    avg: Number(v.avg.toFixed(2)),
    count: v.count,
  }
})

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

  const { payload, label } = props
  const p = payload[0].payload
  return (
    <Paper style={{ padding: 8 }}>
      <strong>{formatDateUnix(label)}</strong> <br />
      {amount(p.avg, 'USD')} <br />
      {p.count} sold
    </Paper>
  )
}

export default function MarketChart() {
  console.log(testdata[0])

  return (
    <div style={{ width: '100%', height: 200 }}>
      <ResponsiveContainer>
        <LineChart data={testdata}>
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
            stroke={primary[800]}
            dot={false}
            strokeWidth={2}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
