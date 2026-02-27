import React, { useContext, useState } from 'react'
import PropTypes from 'prop-types'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Typography from '@mui/material/Typography'
import Divider from '@mui/material/Divider'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import Link from '@/components/Link'
import useLocalStorage from './useLocalStorage'

const targetUpdateID = 1

export default function WhatsNewDialog(props) {
  const { isMobile } = useContext(AppContext)

  const { userID } = props
  const wuid = `whatsnew_id_${userID}`
  const [clientUpdateID, setClientUpdateID] = useLocalStorage(wuid, 0)

  const [open, setOpen] = useState(targetUpdateID > clientUpdateID)
  const handleClose = () => {
    setClientUpdateID(targetUpdateID)
    setOpen(false)
  }

  const handleSubmit = () => {
    handleClose()
  }

  return (
    <Dialog
      fullWidth
      disableEscapeKeyDown
      disableBackdropClick
      fullScreen={isMobile}
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      <DialogTitle id="alert-dialog-title">
        {/* What&apos;s new? */}
        Important announcement
        <DialogCloseButton onClick={handleClose} />
      </DialogTitle>
      <DialogContent>
        <Typography>
          To keep the community fair and a little bit safe for both sellers and buyers. We added
          some ground rules to punish bad behaviors.
          <br />
          <br />
          Please read carefully on{' '}
          <Link
            href="/rules"
            color="secondary"
            target="_blank"
            onClick={() => setClientUpdateID(targetUpdateID)}>
            Rules
          </Link>{' '}
          page.
        </Typography>

        <br />
        <Divider />
        <br />

        <Typography>
          Minor update
          <ul>
            <li>Added time of the latest flagged account on navigation header</li>
          </ul>
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button variant="outlined" color="secondary" onClick={handleSubmit}>
          Got it!
        </Button>
      </DialogActions>
    </Dialog>
  )
}
WhatsNewDialog.propTypes = {
  userID: PropTypes.string,
}
WhatsNewDialog.defaultProps = {
  userID: '',
}
