import React from 'react'
import { Chart } from 'react-charts'
import primary from '@material-ui/core/colors/lightGreen'

const testdata = JSON.parse(
  '[{"avg_sale":6.357142857142857,"day":11,"month":9,"year":2020},{"avg_sale":7.636363636363637,"day":14,"month":9,"year":2020},{"avg_sale":1.5,"day":17,"month":9,"year":2020},{"avg_sale":82,"day":22,"month":9,"year":2020},{"avg_sale":19.285714285714285,"day":2,"month":10,"year":2020},{"avg_sale":12.5,"day":10,"month":10,"year":2020},{"avg_sale":3,"day":13,"month":10,"year":2020},{"avg_sale":7.5,"day":15,"month":10,"year":2020},{"avg_sale":8,"day":21,"month":10,"year":2020},{"avg_sale":14,"day":30,"month":10,"year":2020},{"avg_sale":5.5,"day":2,"month":11,"year":2020},{"avg_sale":15.75,"day":7,"month":11,"year":2020},{"avg_sale":6,"day":8,"month":11,"year":2020},{"avg_sale":9.625,"day":9,"month":11,"year":2020},{"avg_sale":6,"day":12,"month":11,"year":2020},{"avg_sale":34,"day":14,"month":11,"year":2020},{"avg_sale":10,"day":18,"month":11,"year":2020},{"avg_sale":8.5,"day":19,"month":11,"year":2020},{"avg_sale":2.5,"day":21,"month":11,"year":2020},{"avg_sale":35,"day":22,"month":11,"year":2020},{"avg_sale":1.375,"day":24,"month":11,"year":2020},{"avg_sale":3,"day":27,"month":11,"year":2020},{"avg_sale":6.5,"day":30,"month":11,"year":2020}]'
).map(v => {
  return {
    primary: new Date(v.year, v.month, v.day),
    secondary: Number(v.avg_sale.toFixed(2)),
  }
})

export default function MarketChart() {
  console.log(testdata[0])

  const data = React.useMemo(
    () => [
      {
        label: 'Avg Sale',
        data: testdata,
      },
    ],
    []
  )

  const axes = React.useMemo(
    () => [
      { primary: true, type: 'time', position: 'bottom' },
      { position: 'left', type: 'linear' },
    ],
    []
  )

  const getSeriesStyle = React.useCallback(
    () => ({
      color: primary[800],
    }),
    []
  )

  return (
    <div
      style={{
        width: '100%',
        height: 200,
      }}>
      <Chart tooltip dark data={data} axes={axes} getSeriesStyle={getSeriesStyle} />
    </div>
  )
}
