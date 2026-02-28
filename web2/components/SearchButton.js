import React from 'react'
import { makeStyles } from 'tss-react/mui'
import Paper from '@mui/material/Paper'
import InputBase from '@mui/material/InputBase'
import Button from '@mui/material/Button'
import SearchIcon from '@mui/icons-material/Search'
import { styled } from '@mui/material/styles'

const SearchPaper = styled(Paper)(({ theme }) => ({
  padding: '4px 8px 2px',
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
  iconButtons: {
    color: theme.palette.grey[500],
    cursor: 'pointer',
    marginRight: 8,
  },
}))

export default function SearchButton({ style, onClick, ...other }) {
  const { classes } = useStyles()

  return (
    <>
      <SearchPaper
        elevation={0}
        style={style}
        onClick={onClick}
        sx={{
          display: {
            xs: 'none',
            md: 'inherit',
          },
        }}>
        <SearchIcon className={classes.iconButtons} />
        <Input
          size="small"
          placeholder="Search..."
          variant="outlined"
          color="secondary"
          onChange={onClick}
          value=""
          {...other}
        />
      </SearchPaper>
      <Button
        variant="outlined"
        sx={{
          display: {
            width: 36,
            height: 36,
            xs: 'inherit',
            md: 'none',
          },
        }}
        onClick={onClick}>
        <SearchIcon fontSize="small" />
      </Button>
    </>
  )
}
