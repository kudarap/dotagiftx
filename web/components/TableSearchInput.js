import React from 'react'
import PropTypes from 'prop-types'
import debounce from 'lodash/debounce'
import { makeStyles } from '@material-ui/core/styles'
import Paper from '@material-ui/core/Paper'
import InputBase from '@material-ui/core/InputBase'
import CircularProgress from '@material-ui/core/CircularProgress'
import SearchIcon from '@material-ui/icons/Search'
import CloseIcon from '@material-ui/icons/Close'

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

export default function TableSearchInput({ onInput, loading, ...other }) {
  const classes = useStyles()

  const [value, setValue] = React.useState('')

  const debounceSearch = debounce(onInput, 500)
  const handleSearchInput = e => {
    const v = e.target.value
    setValue(v)
    debounceSearch(v)
  }

  const handleSearchClear = () => {
    setValue('')
    debounceSearch('')
  }

  return (
    <Paper className={classes.root} elevation={0}>
      {loading ? (
        <CircularProgress color="secondary" size={24} style={{ marginLeft: 1 }} />
      ) : (
        <SearchIcon className={classes.icon} />
      )}
      <InputBase className={classes.input} value={value} {...other} onInput={handleSearchInput} />

      {value && <CloseIcon style={{ cursor: 'pointer' }} onClick={handleSearchClear} />}
    </Paper>
  )
}
TableSearchInput.propTypes = {
  onInput: PropTypes.func,
  loading: PropTypes.bool,
}
TableSearchInput.defaultProps = {
  onInput: () => {},
  loading: false,
}
