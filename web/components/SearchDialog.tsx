import React, { useContext, useState } from 'react'
import { useRouter } from 'next/router'
import InputBase from '@mui/material/InputBase'
import Typography from '@mui/material/Typography'
import Dialog from '@mui/material/Dialog'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Grid from '@mui/material/Grid'
import Divider from '@mui/material/Divider'
import SearchIcon from '@mui/icons-material/Search'
import useSWR from 'swr'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import Link from './Link'
import { fetcherBase, STATS_TOP_KEYWORDS } from '@/service/api'

const sanitizeInput = (unsafe: string) => {
  if (typeof unsafe !== 'string') return ''

  return unsafe.replace(/[^0-9a-zA-Z- ']/g, '')
}

function SearchDialog({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { isMobile } = useContext(AppContext)

  const { data: topKeywords } = useSWR<{ keyword: string }[]>(STATS_TOP_KEYWORDS, fetcherBase)

  const [keyword, setKeyword] = useState('')
  const router = useRouter()
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onClose()
    router.push(`/search?q=${keyword}`)
  }

  return (
    <Dialog
      fullWidth
      fullScreen={isMobile}
      open={open}
      onClose={onClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title" component="form" onSubmit={handleSubmit}>
        <InputBase
          autoFocus
          fullWidth
          sx={{ fontSize: '1.1em' }}
          startAdornment={<SearchIcon sx={{ mr: 2, fontSize: '1.1em' }} />}
          endAdornment={<DialogCloseButton sx={{ fontSize: '1.1em' }} onClick={onClose} />}
          placeholder="Search for item, hero, or treasure"
          onChange={e => setKeyword(sanitizeInput(e.target.value))}
        />
      </DialogTitle>
      <Divider />
      <DialogContent sx={{ pb: 6 }}>
        <Typography variant="h6" sx={{ mb: 2 }}>
          Top keywords
        </Typography>
        <Grid container spacing={{ xs: 2, sm: 1 }}>
          {topKeywords &&
            topKeywords.map(item => (
              <Grid
                key={item.keyword}
                size={{
                  sm: 6,
                  xs: 12,
                }}
              >
                <Link
                  href={`/search?q=${item.keyword}`}
                  style={{ textTransform: 'capitalize' }}
                  onClick={onClose}
                >
                  {item.keyword}
                </Link>
              </Grid>
            ))}
        </Grid>
      </DialogContent>
    </Dialog>
  )
}

export default SearchDialog
