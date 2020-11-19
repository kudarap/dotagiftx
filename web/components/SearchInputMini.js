import React from 'react'
import PropTypes from 'prop-types'
import querystring from 'querystring'
import has from 'lodash/has'
import { useRouter } from 'next/router'
import { makeStyles } from '@material-ui/core/styles'
import Paper from '@material-ui/core/Paper'
import InputBase from '@material-ui/core/InputBase'
import SearchIcon from '@material-ui/icons/Search'
import CloseIcon from '@material-ui/icons/Close'

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
    padding: '4px 12px 2px',
    marginBottom: 3,
    display: 'flex',
    alignItems: 'center',
    backgroundColor: theme.palette.grey[100],
    // border: `1px solid ${theme.palette.background.paper}`,
    // '&:hover': {
    //   borderColor: theme.palette.grey[700],
    // },
  },
  input: {
    [theme.breakpoints.down('sm')]: {
      height: 42,
    },
    margin: '0 auto',
    color: theme.palette.grey[800],
    // backgroundColor: theme.palette.background.paper,
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

export default function SearchInput({ value, onChange, onSubmit, onClear, ...other }) {
  const classes = useStyles()

  const router = useRouter()
  const { query } = router
  const [keyword, setKeyword] = React.useState(query.q || '')
  React.useEffect(() => {
    if (!has(query, 'q')) {
      return
    }

    setKeyword(query.q)
  }, [query.q])

  const handleChange = ({ target }) => {
    const v = target.value
    setKeyword(v)
    onChange(v)
    // routerPush(v)
  }
  const handleClearValue = () => {
    setKeyword('')
    onChange('')
    router.push(`/search`)
  }
  const handleSubmit = e => {
    e.preventDefault()
    onSubmit(keyword)
    const f = { q: keyword }
    if (query.sort) {
      f.sort = query.sort
    }

    router.push(`/search?${querystring.encode(f)}`)
  }

  return (
    <Paper onSubmit={handleSubmit} className={classes.root} component="form" elevation={0}>
      <InputBase
        onInput={handleChange}
        value={keyword}
        className={classes.input}
        size="small"
        placeholder="Search for item name, hero, treasure"
        variant="outlined"
        color="secondary"
        fullWidth
        {...other}
      />

      {keyword ? (
        <CloseIcon className={classes.iconButtons} onClick={handleClearValue} />
      ) : (
        <SearchIcon className={classes.iconButtons} onClick={handleSubmit} />
      )}
    </Paper>
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
