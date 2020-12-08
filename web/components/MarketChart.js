import React from 'react'
import { Chart } from 'react-charts'
import primary from '@material-ui/core/colors/lightGreen'

const testdata = JSON.parse(
  // '[{"avg":1.5,"day":8,"month":11,"year":2020},{"avg":1.5,"day":24,"month":11,"year":2020},{"avg":2.5,"day":26,"month":11,"year":2020},{"avg":1.5,"day":5,"month":12,"year":2020},{"avg":1.5,"day":7,"month":12,"year":2020}]'
  '[{"avg":13.7,"day":14,"month":9,"year":2020},{"avg":82,"day":14,"month":10,"year":2020},{"avg":7.1,"day":15,"month":10,"year":2020},{"avg":15,"day":25,"month":10,"year":2020},{"avg":10,"day":28,"month":10,"year":2020},{"avg":1.5,"day":3,"month":11,"year":2020},{"avg":10,"day":5,"month":11,"year":2020},{"avg":9,"day":8,"month":11,"year":2020},{"avg":4.25,"day":9,"month":11,"year":2020},{"avg":4.5,"day":13,"month":11,"year":2020},{"avg":10.75,"day":14,"month":11,"year":2020},{"avg":8.666666666666666,"day":18,"month":11,"year":2020},{"avg":25.9,"day":19,"month":11,"year":2020},{"avg":3.25,"day":21,"month":11,"year":2020},{"avg":28.7,"day":22,"month":11,"year":2020},{"avg":1.5,"day":24,"month":11,"year":2020},{"avg":3.375,"day":25,"month":11,"year":2020},{"avg":21.25,"day":26,"month":11,"year":2020},{"avg":2.8333333333333335,"day":27,"month":11,"year":2020},{"avg":8.833333333333334,"day":28,"month":11,"year":2020},{"avg":5.1,"day":29,"month":11,"year":2020},{"avg":4.25,"day":1,"month":12,"year":2020},{"avg":14.25,"day":2,"month":12,"year":2020},{"avg":1.6666666666666667,"day":3,"month":12,"year":2020},{"avg":32,"day":4,"month":12,"year":2020},{"avg":2.125,"day":5,"month":12,"year":2020},{"avg":1.5,"day":7,"month":12,"year":2020}]'
)

export default function MarketChart() {
  const data = React.useMemo(
    () => [
      {
        label: 'Avg Sale(USD)',
        data: testdata.map(v => {
          return {
            primary: new Date(v.year, v.month - 1, v.day),
            secondary: Number(v.avg).toFixed(2),
          }
        }),
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
