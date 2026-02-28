import React from 'react'
import PropTypes from 'prop-types'
import debounce from 'lodash/debounce'
import { makeStyles } from 'tss-react/mui'
import Paper from '@mui/material/Paper'
import InputBase from '@mui/material/InputBase'
import CircularProgress from '@mui/material/CircularProgress'
import SearchIcon from '@mui/icons-material/Search'
import CloseIcon from '@mui/icons-material/Close'

const useStyles = makeStyles()(theme => ({
  root: {
    padding: '3px 12px',
    display: 'flex',
    alignItems: 'center',
    opacity: 0.8,
  },
  input: {
    fontSize: theme.typography.body2.fontSize,
    marginLeft: theme.spacing(1),
  },
  icon: {
    color: theme.palette.grey[500],
  },
}))

export default function ActivitySearchInput({ onInput, loading, ...other }) {
  const { classes } = useStyles()

  const [value, setValue] = React.useState('')

  const debounceSearch = React.useCallback(debounce(onInput, 500), [])

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
        <CircularProgress color="secondary" size={24} style={{ marginRight: 2.5 }} />
      ) : (
        <SearchIcon className={classes.icon} />
      )}
      <InputBase className={classes.input} value={value} {...other} onInput={handleSearchInput} />

      {value && <CloseIcon style={{ cursor: 'pointer' }} onClick={handleSearchClear} />}
    </Paper>
  )
}
ActivitySearchInput.propTypes = {
  onInput: PropTypes.func,
  loading: PropTypes.bool,
}
ActivitySearchInput.defaultProps = {
  onInput: () => {},
  loading: false,
}
