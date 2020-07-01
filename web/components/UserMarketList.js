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
import BuyButton from '@/components/BuyButton'
import RarityTag from '@/components/RarityTag'
import TableHeadCell from '@/components/TableHeadCell'

const useStyles = makeStyles(theme => ({
  seller: {
    display: 'inline-flex',
  },
  link: {
    padding: theme.spacing(2),
  },
}))

const testData = {
  data: [
    {
      id: '10ca6012-454f-4b05-8294-764d95c43e20',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10.02,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:30.133+08:00',
      updated_at: '2020-06-24T03:15:30.133+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
    {
      id: '615da0f7-399f-4dd0-aeaa-58dea38c210f',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:43.357+08:00',
      updated_at: '2020-06-24T03:15:43.357+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
    {
      id: '2d229216-d3ac-4331-9fc3-ab55d79136d3',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:40.18+08:00',
      updated_at: '2020-06-24T03:15:40.18+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
    {
      id: '2ac5fc49-439b-4028-8ec8-6f13721d655d',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10.01,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:29.038+08:00',
      updated_at: '2020-06-24T03:15:29.038+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
    {
      id: 'e629c21f-745f-4e6f-83eb-aa0cb53c1572',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:23.161+08:00',
      updated_at: '2020-06-24T03:15:23.161+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
    {
      id: '442b5e53-b372-4dcf-8970-b74984077976',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:39.255+08:00',
      updated_at: '2020-06-24T03:15:39.255+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
    {
      id: '6c15e201-2e0a-4599-abdd-d6077f4ab622',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10.03,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:31.148+08:00',
      updated_at: '2020-06-24T03:15:31.148+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
    {
      id: 'fe69e022-e5d1-4c62-a603-e47a15676d98',
      user_id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
      item_id: '20bab823-2796-4948-879e-518783bb8cff',
      price: 10.04,
      currency: 'USD',
      notes: '',
      status: 200,
      created_at: '2020-06-24T03:15:32.488+08:00',
      updated_at: '2020-06-24T03:15:32.488+08:00',
      user: {
        id: 'f0ee7f4e-aae2-45a2-a599-3dd5f58bbc17',
        steam_id: '76561198088587178',
        name: 'kudarap',
        url: 'https://steamcommunity.com/id/kudarap/',
        avatar: '2639629cc8fe0078393c27848a6f511247ec0195.jpg',
        created_at: '2020-06-18T13:13:43.926+08:00',
        updated_at: '2020-06-18T13:13:43.926+08:00',
      },
      item: {
        id: '20bab823-2796-4948-879e-518783bb8cff',
        slug: 'pitmouse-fraternity-meepo',
        name: 'Pitmouse Fraternity',
        hero: 'Meepo',
        image: '',
        origin: "Collector's Cache 2 2018",
        rarity: 'regular',
        created_at: '2020-06-22T19:32:51.846+08:00',
        updated_at: '2020-06-22T19:32:51.846+08:00',
      },
    },
  ],
  result_count: 8,
  total_count: 8,
}

export default function SimpleTable() {
  const classes = useStyles()

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableHeadCell>Seller Listings</TableHeadCell>
            <TableHeadCell align="right">Price</TableHeadCell>
            <TableHeadCell align="right" />
          </TableRow>
        </TableHead>
        <TableBody>
          {testData.data.map(market => (
            <TableRow key={market.id} hover>
              <TableCell component="th" scope="row" padding="none">
                <Link href="/item/[slug]" as={`/item/${market.item.slug}`} disableUnderline>
                  <div className={classes.link}>
                    <strong>{market.item.name}</strong>
                    <br />
                    <Typography variant="caption" color="textSecondary">
                      {market.item.hero}
                    </Typography>
                    <RarityTag rarity={market.item.rarity} />
                  </div>
                </Link>
              </TableCell>
              <TableCell align="right">
                <Typography variant="body2">${market.price.toFixed(2)}</Typography>
              </TableCell>
              <TableCell align="right">
                <BuyButton variant="contained">Contact Seller</BuyButton>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
