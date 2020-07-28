import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Paper from '@material-ui/core/Paper'
import TextField from '@material-ui/core/TextField'
import Typography from '@material-ui/core/Typography'
import { catalog } from '@/service/api'
import * as format from '@/lib/format'
import Button from '@/components/Button'
import ItemAutoComplete from '@/components/ItemAutoComplete'
import ItemImage from '@/components/ItemImage'

const useStyles = makeStyles(theme => ({
  root: {
    maxWidth: theme.breakpoints.values.sm,
    margin: '0 auto',
    padding: theme.spacing(2),
  },
  itemImage: {
    width: 150,
    height: 100,
    float: 'left',
    marginRight: theme.spacing(1),
  },
}))

export default function MarketForm() {
  const classes = useStyles()

  const [item, setItem] = React.useState({ id: '' })

  const handleItemSelect = val => {
    setItem(val)
    // Get item starting price
    if (val.slug) {
      catalog(val.slug).then(res => setItem(res))
    }
  }

  return (
    <Paper component="form" className={classes.root}>
      <Typography variant="h5" component="h1">
        Selling your item
      </Typography>
      <br />

      <ItemAutoComplete onSelect={handleItemSelect} />
      {/* <TextField */}
      {/*  variant="outlined" */}
      {/*  fullWidth */}
      {/*  required */}
      {/*  color="secondary" */}
      {/*  label="Item name" */}
      {/*  helperText="Search item you want to post from your inventory." */}
      {/*  autoFocus */}
      {/* /> */}
      <br />

      {/* Selected item preview */}
      {item.id && (
        <div>
          <ItemImage
            className={classes.itemImage}
            image={`/300x170/${item.image}`}
            rarity={item.rarity}
            title={item.name}
          />
          <Typography variant="body2" color="textSecondary">
            Origin:{' '}
            <Typography variant="body2" color="textPrimary" component="span">
              {item.origin}
            </Typography>
          </Typography>
          <Typography variant="body2" color="textSecondary">
            Rarity:{' '}
            <Typography variant="body2" color="textPrimary" component="span">
              {item.rarity}
            </Typography>
          </Typography>
          <Typography variant="body2" color="textSecondary">
            Hero:{' '}
            <Typography variant="body2" color="textPrimary" component="span">
              {item.hero}
            </Typography>
          </Typography>
          <Typography variant="body2" color="textSecondary">
            Starting at:{' '}
            <Typography variant="body2" color="textPrimary" component="span">
              {format.amount(item.lowest_ask, 'USD')}
            </Typography>
          </Typography>
          <br />
          <br />
        </div>
      )}

      <div>
        <TextField
          variant="outlined"
          required
          color="secondary"
          label="Price"
          placeholder="1.00"
          type="number"
          helperText="Price value will be on USD."
          style={{ width: '69%' }}
        />
        <TextField
          variant="outlined"
          color="secondary"
          label="Qty"
          type="number"
          defaultValue="1"
          style={{ width: '30%', marginLeft: '1%' }}
        />
      </div>
      <br />
      <TextField
        variant="outlined"
        fullWidth
        color="secondary"
        label="Notes"
        helperText="Keep it short, maximum of 100 characters."
      />
      <br />
      <br />

      <Button variant="contained" fullWidth type="submit" size="large">
        Post Item
      </Button>
    </Paper>
  )
}
