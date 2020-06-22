import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import Link from '@/components/Link'
import TableActions from '@/components/TableActions'
import RarityTag from '@/components/RarityTag'

const useStyles = makeStyles({
  table: {
    // minWidth: 650,
  },
  th: {
    cursor: 'pointer',
  },
  pagination: {
    textAlign: 'right',
  },
})

// background: linear-gradient(#f9ffbf 10%, #fff 90%);
// text-shadow: 0px 0px 10px yellowgreen;
// -webkit-background-clip: text;
// -webkit-text-fill-color: transparent;

// background: linear-gradient(#fdd08e 10%, #fff 90%);
// text-shadow: 0px 0px 10px darkorange;
// -webkit-background-clip: text;
// -webkit-text-fill-color: transparent;

// background: linear-gradient(#F8E8B9 10%, #fff 90%);
// text-shadow: 0px 0px 10px goldenrod;
// -webkit-background-clip: text;
// -webkit-text-fill-color: transparent;
const rarityColorMap = {
  regular: null,
  rare: 'yellowgreen',
  'very rare': 'darkorange',
  'ultra rare': 'goldenrod',
}
const getRarityColor = value => {
  if (value === '') {
    return null
  }

  return rarityColorMap[value]
}

export default function SimpleTable({
  result = {
    data: [],
    result_count: 0,
    total_count: 0,
  },
}) {
  const classes = useStyles()

  return (
    <>
      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="items table">
          <TableHead>
            <TableRow>
              <TableCell>Item Name</TableCell>
              <TableCell align="right">Qty</TableCell>
              <TableCell align="right">Price</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {result.data.map(item => (
              <TableRow key={item.id} hover>
                <TableCell className={classes.th} component="th" scope="row">
                  <Link href="/item/[slug]" as={`/item/${item.slug}`} disableUnderline>
                    <>
                      <strong>{item.name}</strong>
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {item.hero}
                      </Typography>
                      <RarityTag rarity={item.rarity} />
                    </>
                  </Link>
                </TableCell>
                <TableCell align="right">{item.name.length}</TableCell>
                <TableCell align="right">
                  <Typography variant="body2" color="secondary">
                    ${item.hero.length.toFixed(2)}
                  </Typography>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TableActions
        className={classes.pagination}
        colSpan={3}
        count={result.total_count}
        page={1}
      />
    </>
  )
}
