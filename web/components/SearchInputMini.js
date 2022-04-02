import React from 'react'
import PropTypes from 'prop-types'
import querystring from 'querystring'
import has from 'lodash/has'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import Paper from '@mui/material/Paper'
import InputBase from '@mui/material/InputBase'
import SearchIcon from '@mui/icons-material/Search'
import CloseIcon from '@mui/icons-material/Close'

import { styled } from '@mui/material/styles'

const SearchPaper = styled(Paper)(({ theme }) => ({
  padding: '4px 8px 2px 12px',
  marginBottom: 3,
  display: 'flex',
  alignItems: 'center',
  backgroundColor: theme.palette.background.default,
  width: 325,
}))

const Input = styled(InputBase)(({ theme }) => ({
  [theme.breakpoints.down('md')]: {
    height: 39,
  },
  margin: '0 auto',
  color: theme.palette.grey[100],
}))

const useStyles = makeStyles()(theme => ({
  actionIcons: {
    color: theme.palette.grey[500],
  },
  iconButtons: {
    color: theme.palette.grey[500],
    cursor: 'pointer',
  },
}))

export default function SearchInput({ value, onChange, onSubmit, onClear, style, ...other }) {
  const { classes } = useStyles()

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
    <SearchPaper
      onSubmit={handleSubmit}
      className={classes.root}
      component="form"
      elevation={0}
      style={style}>
      <Input
        onInput={handleChange}
        value={keyword}
        className={classes.input}
        size="small"
        placeholder="Search..."
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
    </SearchPaper>
  )
}
SearchInput.propTypes = {
  value: PropTypes.string,
  onChange: PropTypes.func,
  onSubmit: PropTypes.func,
  onClear: PropTypes.func,
  style: PropTypes.object,
}
SearchInput.defaultProps = {
  value: '',
  onChange: () => {},
  onSubmit: () => {},
  onClear: () => {},
  style: {},
}
