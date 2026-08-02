import React from 'react'
import debounce from 'lodash/debounce'
import { makeStyles } from 'tss-react/mui'
import Paper from '@mui/material/Paper'
import InputBase from '@mui/material/InputBase'
import type { InputBaseProps } from '@mui/material/InputBase'
import CircularProgress from '@mui/material/CircularProgress'
import SearchIcon from '@mui/icons-material/Search'
import CloseIcon from '@mui/icons-material/Close'

const useStyles = makeStyles()(theme => ({
  root: {
    padding: '3px 12px',
    display: 'flex',
    alignItems: 'center',
    backgroundColor: theme.palette.background.default,
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

interface TableSearchInputProps extends Omit<InputBaseProps, 'onInput'> {
  value?: string
  onInput?: (value: string) => void
  loading?: boolean
}

export default function TableSearchInput({ value: externalValue, onInput, loading, ...other }: TableSearchInputProps) {
  const { classes } = useStyles()

  const [value, setValue] = React.useState(externalValue || '')

  const debounceSearch = React.useMemo(() => debounce(onInput || (() => {}), 500), [onInput])

  const handleSearchInput = (e: { target: { value: string } }) => {
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
      <InputBase
        className={classes.input}
        value={value}
        {...other}
        onInput={handleSearchInput as unknown as NonNullable<InputBaseProps['onInput']>}
      />

      {value && <CloseIcon style={{ cursor: 'pointer' }} onClick={handleSearchClear} />}
    </Paper>
  )
}
