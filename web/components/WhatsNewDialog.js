import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import RemoveIcon from '@material-ui/icons/Close'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'

export default function WhatsNewDialog(props) {
  const { isMobile } = useContext(AppContext)

  const { onClose } = props
  const handleClose = () => {
    onClose()
  }

  const handleSubmit = () => {}

  const { open } = props
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
        Keeping community a bit less abusive
        <DialogCloseButton onClick={handleClose} />
      </DialogTitle>
      <DialogContent>
        <Typography>Seller</Typography>
      </DialogContent>
      <DialogActions>
        <Button variant="outlined" color="secondary" onClick={handleSubmit}>
          Done
        </Button>
      </DialogActions>
    </Dialog>
  )
}
WhatsNewDialog.propTypes = {
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
WhatsNewDialog.defaultProps = {
  open: false,
  onClose: () => {},
}
