import React, { useContext, useState } from 'react'
import PropTypes from 'prop-types'
import { useRouter } from 'next/router'
import InputBase from '@mui/material/InputBase'
import Typography from '@mui/material/Typography'
import Dialog from '@mui/material/Dialog'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Grid from '@mui/material/Grid'
import Divider from '@mui/material/Divider'
import SearchIcon from '@mui/icons-material/Search'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import Link from './Link'

const tempTopKeywords = [
  {
    keyword: 'nemestice',
    score: 262,
  },
  {
    keyword: 'pudge',
    score: 149,
  },
  {
    keyword: 'sven',
    score: 139,
  },
  {
    keyword: 'tiny',
    score: 127,
  },
  {
    keyword: '2015',
    score: 106,
  },
  {
    keyword: 'immortal',
    score: 105,
  },
  {
    keyword: 'night stalker',
    score: 86,
  },
  {
    keyword: 'shadow fiend',
    score: 86,
  },
  {
    keyword: 'ember',
    score: 85,
  },
  {
    keyword: 'huskar',
    score: 85,
  },
  {
    keyword: 'undying',
    score: 78,
  },
  {
    keyword: 'Pudge',
    score: 77,
  },
]

function SearchDialog({ open, onClose }) {
  const { isMobile } = useContext(AppContext)

  const [keyword, setKeyword] = useState('')
  const router = useRouter()
  const handleSubmit = e => {
    e.preventDefault()
    onClose()
    router.push(`/search?q=${keyword}`)
  }

  return (
    <>
      <Dialog
        fullWidth
        fullScreen={isMobile}
        open={open}
        onClose={onClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
        <DialogTitle id="alert-dialog-title" component="form" onSubmit={handleSubmit}>
          <InputBase
            autoFocus
            fullWidth
            sx={{ fontSize: '1.1em' }}
            startAdornment={<SearchIcon sx={{ mr: 2, fontSize: '1.1em' }} />}
            endAdornment={<DialogCloseButton sx={{ fontSize: '1.1em' }} onClick={onClose} />}
            placeholder="Search..."
            onChange={e => setKeyword(e.target.value)}
          />
        </DialogTitle>
        <Divider />
        <DialogContent sx={{ pb: 6 }}>
          <Typography variant="h6" sx={{ mb: 2 }}>
            Top keywords
          </Typography>
          <Grid container spacing={{ xs: 2, sm: 1 }}>
            {tempTopKeywords.map(item => (
              <Grid key={item.keyword} item sm={6} xs={12}>
                <Link
                  href={`/search?q=${item.keyword}`}
                  style={{ textTransform: 'capitalize' }}
                  onClick={onClose}>
                  {item.keyword}
                </Link>
              </Grid>
            ))}
          </Grid>
        </DialogContent>
      </Dialog>
    </>
  )
}
SearchDialog.propTypes = {
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
SearchDialog.defaultProps = {
  open: false,
  onClose: () => {},
}

export default SearchDialog
