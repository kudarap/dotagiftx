import React, { useContext, useState } from 'react'
import PropTypes from 'prop-types'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Grid from '@mui/material/Grid'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemText from '@mui/material/ListItemText'
import Divider from '@mui/material/Divider'
import SearchIcon from '@mui/icons-material/Search'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import SearchInput from '@/components/SearchInput'
import Link from '@/components/Link'
import { InputBase } from '@mui/material'

function generate(element) {
  return [0, 1, 2, 3, 4].map(value =>
    React.cloneElement(element, {
      key: value,
    })
  )
}

export default function SearchDialog() {
  const { isMobile } = useContext(AppContext)

  const [open, setOpen] = useState(false)
  const handleClose = () => setOpen(false)

  const handleSubmit = () => {
    handleClose()
  }

  return (
    <>
      <Button variant="outlined" sx={{ px: 1, py: 0.8, minWidth: 0 }} onClick={() => setOpen(true)}>
        <SearchIcon fontSize="small" />
      </Button>

      <Dialog
        fullWidth
        fullScreen={isMobile}
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
        <DialogTitle id="alert-dialog-title">
          <InputBase
            autoFocus
            sx={{ fontSize: '1.1em' }}
            startAdornment={<SearchIcon sx={{ mr: 2, fontSize: '1.1em' }} />}
            endAdornment={<DialogCloseButton sx={{ fontSize: '1.1em' }} onClick={handleClose} />}
            placeholder="Search..."
            fullWidth
          />
        </DialogTitle>
        <Divider />
        <DialogContent>
          <Grid container spacing={2}>
            <Grid item xs={6}>
              <List>
                {generate(
                  <ListItem>
                    <ListItemText primary="Single-line item" />
                  </ListItem>
                )}
              </List>
            </Grid>
          </Grid>
        </DialogContent>
      </Dialog>
    </>
  )
}
SearchDialog.propTypes = {
  userID: PropTypes.string,
}
SearchDialog.defaultProps = {
  userID: '',
}
