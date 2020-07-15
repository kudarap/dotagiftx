import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'
import InputAdornment from '@material-ui/core/InputAdornment'
import SearchIcon from '@material-ui/icons/Search'
import CloseIcon from '@material-ui/icons/Close'

const useStyles = makeStyles(theme => ({
  main: {
    // marginTop: theme.spacing(4),
  },
  searchBar: {
    margin: '0 auto',
    '& .MuiInputBase-root': {
      color: theme.palette.grey[300],
      backgroundColor: theme.palette.background.paper,
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

export default function SearchInput(props) {
  const { value, onChange, onSubmit, onClear, ...other } = props

  const classes = useStyles()

  const [keyword, setKeyword] = React.useState(value)

  const handleChange = ({ target }) => {
    const v = target.value
    setKeyword(v)
    onChange(v)
  }
  const handleClearValue = () => {
    setKeyword('')
    onChange('')
    onClear()
  }
  const handleSubmit = e => {
    e.preventDefault()
    onSubmit(keyword)
  }

  return (
    <form onSubmit={handleSubmit} style={{ flexGrow: 1 }}>
      <TextField
        size="small"
        className={classes.searchBar}
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              {keyword !== '' ? (
                <CloseIcon className={classes.iconButtons} onClick={handleClearValue} />
              ) : (
                <SearchIcon className={classes.iconButtons} onClick={handleSubmit} />
              )}
            </InputAdornment>
          ),
        }}
        placeholder="Search for item name, hero, origin"
        variant="outlined"
        color="secondary"
        fullWidth
        value={keyword}
        onInput={handleChange}
        {...other}
      />
    </form>
  )
}
SearchInput.propTypes = {
  value: PropTypes.string,
  onChange: PropTypes.func,
  onSubmit: PropTypes.func,
  onClear: PropTypes.func,
}
SearchInput.defaultProps = {
  value: '',
  onChange: () => {},
  onSubmit: () => {},
  onClear: () => {},
}
