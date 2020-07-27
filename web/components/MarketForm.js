import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Paper from '@material-ui/core/Paper'
import TextField from '@material-ui/core/TextField'
import Typography from '@material-ui/core/Typography'
import Button from '@/components/Button'
import ItemAutoComplete from '@/components/ItemAutoComplete'

const useStyles = makeStyles(theme => ({
  root: {
    maxWidth: theme.breakpoints.values.sm,
    margin: '0 auto',
    padding: theme.spacing(2),
  },
}))

export default function MarketForm() {
  const classes = useStyles()

  return (
    <Paper component="form" className={classes.root}>
      <Typography variant="h5" component="h1">
        Selling your item
      </Typography>
      <br />

      <ItemAutoComplete />
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
