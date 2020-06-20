import React from 'react'
import { makeStyles, useTheme } from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'
import InputAdornment from '@material-ui/core/InputAdornment'
import SearchIcon from '@material-ui/icons/Search'
import ArrowForwardIcon from '@material-ui/icons/ArrowForward'
import CloseIcon from '@material-ui/icons/Close'

const useStyles = makeStyles(theme => ({
  main: {
    marginTop: theme.spacing(4),
  },
  searchBar: {
    margin: '0 auto',
    marginBottom: theme.spacing(4),
  },
  verticalDivider: {
    borderRight: `1px solid ${theme.palette.divider}`,
    height: 40,
    margin: theme.spacing(0, 1.5),
  },
  actionIcons: {
    color: theme.palette.text.hint,
  },
  iconButtons: {
    color: theme.palette.text.hint,
    cursor: 'pointer',
  },
}))

export default function SearchInput() {
  const classes = useStyles()

  const [value, setValue] = React.useState('')

  const handleChange = e => {
    setValue(e.target.value)
  }

  return (
    <>
      <TextField
        className={classes.searchBar}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon className={classes.actionIcons} />
            </InputAdornment>
          ),
          endAdornment: (
            <InputAdornment position="end">
              {value !== '' && (
                <>
                  <CloseIcon className={classes.iconButtons} />
                  <span className={classes.verticalDivider} />
                </>
              )}

              <ArrowForwardIcon className={classes.iconButtons} />
            </InputAdornment>
          ),
        }}
        placeholder="Search Item, Hero, Treasure..."
        helperText="Search on 100+ for sale items"
        variant="outlined"
        color="secondary"
        fullWidth
        autoFocus
        value={value}
        onChange={handleChange}
      />
    </>
  )
}
