import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from 'tss-react/mui'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'
import InputAdornment from '@mui/material/InputAdornment'
import SearchIcon from '@mui/icons-material/Search'
import ArrowForwardIcon from '@mui/icons-material/ArrowForward'
import CloseIcon from '@mui/icons-material/Close'

const useStyles = makeStyles()(theme => ({
  main: {
    // marginTop: theme.spacing(4),
  },
  searchBar: {
    margin: '0 auto',
    // marginBottom: theme.spacing(4),
    '& .MuiInputBase-root': {
      color: theme.palette.grey[800],
      backgroundColor: theme.palette.common.white,
      // backgroundColor: theme.palette.app.white,
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
  label: {
    margin: theme.spacing(0.4, 0, 0, 1.8),
    display: 'block',
  },
}))

export default function SearchInput(props) {
  const { value, onChange, onSubmit, onClear, label, ...other } = props

  const { classes } = useStyles()

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
    <form onSubmit={handleSubmit}>
      <TextField
        id="search_input"
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
        placeholder="Search for item name, hero, treasure"
        variant="outlined"
        color="secondary"
        fullWidth
        value={keyword}
        onChange={handleChange}
        {...other}
      />
      {label && (
        <Typography
          className={classes.label}
          variant="caption"
          color="textSecondary"
          component="label"
          htmlFor="search_input">
          {label}
        </Typography>
      )}
    </form>
  )
}
SearchInput.propTypes = {
  value: PropTypes.string,
  label: PropTypes.string,
  onChange: PropTypes.func,
  onSubmit: PropTypes.func,
  onClear: PropTypes.func,
}
SearchInput.defaultProps = {
  value: '',
  label: 'Search Giftable items',
  onChange: () => {},
  onSubmit: () => {},
  onClear: () => {},
}
