import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import moment from 'moment'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import Link from '@/components/Link'
import RarityTag from '@/components/RarityTag'

const useStyles = makeStyles({
  table: {
    // minWidth: 650,
  },
})

const testData = {
  data: [
    {
      id: '88f2fcdc-5f04-4afe-bf01-a28bb5addf21',
      slug: 'dread-compact-warlock',
      name: 'Dread Compact',
      hero: 'Warlock',
      image: '',
      origin: "Collector's Cache 2018",
      rarity: 'rare',
      created_at: '2020-06-22T17:09:29.983+08:00',
      updated_at: '2020-06-22T17:09:29.983+08:00',
    },
    {
      id: '29ebdb2f-859b-4822-8f00-883a09319f4f',
      slug: 'visions-of-the-lifted-veil-phantom-assassin',
      name: 'Visions of the Lifted Veil',
      hero: 'Phantom Assassin',
      image: '',
      origin: "Collector's Cache 2018",
      rarity: 'regular',
      created_at: '2020-06-22T17:09:17.113+08:00',
      updated_at: '2020-06-22T17:09:17.113+08:00',
    },
    {
      id: '8c282a8b-eb0b-4068-b578-d4cb7c0d9eff',
      slug: 'grasp-of-the-riven-exile-weaver',
      name: 'Grasp of the Riven Exile',
      hero: 'Weaver',
      image: '',
      origin: "Collector's Cache 2018",
      rarity: 'regular',
      created_at: '2020-06-22T17:09:07.719+08:00',
      updated_at: '2020-06-22T17:09:07.719+08:00',
    },
    {
      id: 'a13a01b4-1c3a-489e-853b-6b07005fee7c',
      slug: 'fate-meridian-invoker',
      name: 'Fate Meridian',
      hero: 'Invoker',
      image: '',
      origin: "Collector's Cache 2018",
      rarity: 'regular',
      created_at: '2020-06-22T17:08:59.153+08:00',
      updated_at: '2020-06-22T17:08:59.153+08:00',
    },
    {
      id: '2ec6d2e3-2101-46e9-9aa6-47d7e7d010e2',
      slug: 'raptures-of-the-abyssal-kin-queen-of-paing',
      name: 'Raptures of the Abyssal Kin',
      hero: 'Queen of Paing',
      image: '',
      origin: "Collector's Cache 2018",
      rarity: 'regular',
      created_at: '2020-06-22T17:08:51.023+08:00',
      updated_at: '2020-06-22T17:08:51.023+08:00',
    },
  ],
  result_count: 5,
  total_count: 39,
}

export default function SimpleTable() {
  const classes = useStyles()

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Item Name</TableCell>
            <TableCell align="right">Date</TableCell>
            <TableCell align="right">Price</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {testData.data.map(item => (
            <TableRow key={item.id} hover>
              <TableCell component="th" scope="row">
                <Link href="/item/[slug]" as={`/item/${item.slug}`} disableUnderline>
                  <>
                    <strong>{item.name}</strong>
                    <RarityTag rarity={item.rarity} />
                    <br />
                    <Typography variant="caption" color="textSecondary">
                      {item.hero}
                    </Typography>
                  </>
                </Link>
              </TableCell>
              <TableCell align="right">{moment(item.created_at).fromNow()}</TableCell>
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
