import React from 'react'
import PropTypes from 'prop-types'
import Router from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
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
    '& .MuiInputBase-root': {
      color: theme.palette.grey[800],
      backgroundColor: theme.palette.app.white,
    },
  },
  verticalDivider: {
    borderRight: `1px solid ${theme.palette.grey[300]}`,
    height: 40,
    margin: theme.spacing(0, 1.5),
  },
  actionIcons: {
    color: theme.palette.grey[500],
  },
  iconButtons: {
    color: theme.palette.grey[500],
    cursor: 'pointer',
  },
}))

export default function SearchInput({ value }) {
  const classes = useStyles()

  const [keyword, setKeyword] = React.useState(value)

  const handleChange = e => {
    setKeyword(e.target.value)
  }
  const handleClearValue = () => {
    setKeyword('')
  }
  const handleSubmit = e => {
    e.preventDefault()
    Router.push(`/search?q=${keyword}`)
  }

  return (
    <form onSubmit={handleSubmit}>
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
              {keyword !== '' && (
                <>
                  <CloseIcon className={classes.iconButtons} onClick={handleClearValue} />
                  <span className={classes.verticalDivider} />
                </>
              )}

              <ArrowForwardIcon className={classes.iconButtons} onClick={handleSubmit} />
            </InputAdornment>
          ),
        }}
        placeholder="Search Item, Hero, Treasure..."
        helperText="Search on 100+ for sale items"
        variant="outlined"
        color="secondary"
        fullWidth
        autoFocus
        value={keyword}
        onChange={handleChange}
      />
    </form>
  )
}
SearchInput.propTypes = {
  value: PropTypes.string,
}
SearchInput.defaultProps = {
  value: '',
}
