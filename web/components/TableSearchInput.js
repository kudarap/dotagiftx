import React from 'react'
import PropTypes from 'prop-types'
import debounce from 'lodash/debounce'
import { makeStyles } from '@material-ui/core/styles'
import Paper from '@material-ui/core/Paper'
import InputBase from '@material-ui/core/InputBase'
import SearchIcon from '@material-ui/icons/Search'

const useStyles = makeStyles(theme => ({
  root: {
    padding: '3px 12px',
    display: 'flex',
    alignItems: 'center',
    backgroundColor: theme.palette.primary.dark,
    opacity: 0.8,
    margin: theme.spacing(1),
  },
  input: {
    fontSize: theme.typography.body2.fontSize,
    marginLeft: theme.spacing(1),
  },
  icon: {
    color: theme.palette.grey[500],
  },
}))

export default function TableSearchInput({ onInput, ...other }) {
  const classes = useStyles()

  const debounceSearch = debounce(onInput, 500)
  const handleSearchInput = e => {
    debounceSearch(e.target.value)
  }

  return (
    <Paper className={classes.root} elevation={0}>
      <SearchIcon className={classes.icon} />
      <InputBase className={classes.input} {...other} onInput={handleSearchInput} />
    </Paper>
  )
}
TableSearchInput.propTypes = {
  onInput: PropTypes.func,
}
TableSearchInput.defaultProps = {
  onInput: () => {},
}
