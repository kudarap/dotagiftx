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

const useStyles = makeStyles({
  table: {
    // minWidth: 650,
  },
})

const testData = {
  data: [
    {
      id: 'f5d0211d-3377-457e-a742-afa5ea16a9bf',
      slug: 'pipe-of-dezun-dazzle',
      name: 'Pipe of Dezun',
      hero: 'Dazzle',
      image: '',
      origin: 'Treasure of the Cryptic Beacon',
      rarity: 'very rare',
      created_at: '2020-06-22T17:19:57.58+08:00',
      updated_at: '2020-06-22T17:19:57.58+08:00',
    },
    {
      id: '5bde1aca-2300-4e16-b62a-772fa55c6d1a',
      slug: 'diabolic-aspect-chaos-knight',
      name: 'Diabolic Aspect',
      hero: 'Chaos Knight',
      image: '',
      origin: 'Treasure of the Cryptic Beacon',
      rarity: 'regular',
      created_at: '2020-06-22T17:19:33.978+08:00',
      updated_at: '2020-06-22T17:19:33.978+08:00',
    },
    {
      id: '4faf069d-884b-4c97-bc4b-3e26eee3d89e',
      slug: 'shattered-greatsword-sven',
      name: 'Shattered Greatsword',
      hero: 'Sven',
      image: '',
      origin: 'Treasure of the Cryptic Beacon',
      rarity: 'regular',
      created_at: '2020-06-22T17:19:07.769+08:00',
      updated_at: '2020-06-22T17:19:07.769+08:00',
    },
    {
      id: '01ed7e48-ff44-4428-b1f8-cbb3319bbb67',
      slug: 'vespidun-hunter-killer-gyrocopter',
      name: 'Vespidun Hunter-Killer',
      hero: 'Gyrocopter',
      image: '',
      origin: "Collector's Cache 2018",
      rarity: 'ultra rare',
      created_at: '2020-06-22T17:10:04.586+08:00',
      updated_at: '2020-06-22T17:10:04.586+08:00',
    },
    {
      id: '9c0c10ad-a642-466a-be22-8eb8bde5e71b',
      slug: 'endowments-of-the-lucent-canopy-shadow-shaman',
      name: 'Endowments of the Lucent Canopy',
      hero: 'Shadow Shaman',
      image: '',
      origin: "Collector's Cache 2018",
      rarity: 'very rare',
      created_at: '2020-06-22T17:09:43.25+08:00',
      updated_at: '2020-06-22T17:09:43.25+08:00',
    },
  ],
  result_count: 5,
  total_count: 39,
}

export default function SimpleTable() {
  const classes = useStyles()

  return (
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
          {testData.data.slice(0, 5).map(item => (
            <TableRow key={item.id}>
              <TableCell component="th" scope="row">
                <Link href="/item/[slug]" as={`/item/${item.slug}`} disableUnderline>
                  <>
                    <strong>{item.name}</strong>
                    <br />
                    <Typography variant="caption" color="textSecondary">
                      {item.hero}
                    </Typography>
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
  )
}
