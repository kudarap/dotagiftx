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

const useStyles = makeStyles({
  table: {
    // minWidth: 650,
  },
})

const testData = {
  data: [
    {
      id: 'fa4757c2-c5e0-4ced-a67a-7a96fe533952',
      slug: 'allure-of-the-faeshade-flower-dark-willow',
      name: 'Allure of the Faeshade Flower',
      hero: 'Dark Willow',
      image: '',
      origin: "Collector's Cache 2018",
      created_at: '2020-06-19T00:59:58.051+08:00',
      updated_at: '2020-06-19T00:59:58.051+08:00',
    },
    {
      id: 'f8b117bb-3a1f-469b-b2c8-8014df13ea09',
      slug: 'defender-of-the-ruins-disruptor',
      name: 'Defender of the Ruins',
      hero: 'Disruptor',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:29.685+08:00',
      updated_at: '2020-06-19T18:06:29.685+08:00',
    },
    {
      id: '09175f2f-78f8-4abe-8bdd-2ed46be359a4',
      slug: 'gothink-whisper-phantom-assassin',
      name: 'Gothink Whisper',
      hero: 'Phantom Assassin',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:02.118+08:00',
      updated_at: '2020-06-19T18:06:02.118+08:00',
    },
    {
      id: 'fa4757c2-c5e0-4ced-a67a-7a96fe533952',
      slug: 'allure-of-the-faeshade-flower-dark-willow',
      name: 'Allure of the Faeshade Flower',
      hero: 'Dark Willow',
      image: '',
      origin: "Collector's Cache 2018",
      created_at: '2020-06-19T00:59:58.051+08:00',
      updated_at: '2020-06-19T00:59:58.051+08:00',
    },
    {
      id: 'f8b117bb-3a1f-469b-b2c8-8014df13ea09',
      slug: 'defender-of-the-ruins-disruptor',
      name: 'Defender of the Ruins',
      hero: 'Disruptor',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:29.685+08:00',
      updated_at: '2020-06-19T18:06:29.685+08:00',
    },
    {
      id: '09175f2f-78f8-4abe-8bdd-2ed46be359a4',
      slug: 'gothink-whisper-phantom-assassin',
      name: 'Gothink Whisper',
      hero: 'Phantom Assassin',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:02.118+08:00',
      updated_at: '2020-06-19T18:06:02.118+08:00',
    },
    {
      id: 'f8b117bb-3a1f-469b-b2c8-8014df13ea09',
      slug: 'defender-of-the-ruins-disruptor',
      name: 'Defender of the Ruins',
      hero: 'Disruptor',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:29.685+08:00',
      updated_at: '2020-06-19T18:06:29.685+08:00',
    },
    {
      id: '09175f2f-78f8-4abe-8bdd-2ed46be359a4',
      slug: 'gothink-whisper-phantom-assassin',
      name: 'Gothink Whisper',
      hero: 'Phantom Assassin',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:02.118+08:00',
      updated_at: '2020-06-19T18:06:02.118+08:00',
    },
    {
      id: 'fa4757c2-c5e0-4ced-a67a-7a96fe533952',
      slug: 'allure-of-the-faeshade-flower-dark-willow',
      name: 'Allure of the Faeshade Flower',
      hero: 'Dark Willow',
      image: '',
      origin: "Collector's Cache 2018",
      created_at: '2020-06-19T00:59:58.051+08:00',
      updated_at: '2020-06-19T00:59:58.051+08:00',
    },
    {
      id: 'f8b117bb-3a1f-469b-b2c8-8014df13ea09',
      slug: 'defender-of-the-ruins-disruptor',
      name: 'Defender of the Ruins',
      hero: 'Disruptor',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:29.685+08:00',
      updated_at: '2020-06-19T18:06:29.685+08:00',
    },
    {
      id: '09175f2f-78f8-4abe-8bdd-2ed46be359a4',
      slug: 'gothink-whisper-phantom-assassin',
      name: 'Gothink Whisper',
      hero: 'Phantom Assassin',
      image: '',
      origin: "Collector's Cache 2019",
      created_at: '2020-06-19T18:06:02.118+08:00',
      updated_at: '2020-06-19T18:06:02.118+08:00',
    },
  ],
  result_count: 10,
  total_count: 332,
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
          {testData.data.slice(0, 5).map(item => (
            <TableRow key={item.id}>
              <TableCell component="th" scope="row">
                <strong>{item.name}</strong>
                <br />
                <Typography variant="caption" color="textSecondary">
                  {item.hero}
                </Typography>
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
