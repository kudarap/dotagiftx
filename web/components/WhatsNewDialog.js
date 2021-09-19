import React, { useContext, useState } from 'react'
import PropTypes from 'prop-types'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import Divider from '@material-ui/core/Divider'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import Link from '@/components/Link'

const targetUpdateID = 1

function useLocalStorage(key, initialValue) {
  // State to store our value
  // Pass initial state function to useState so logic is only executed once
  const [storedValue, setStoredValue] = useState(() => {
    if (typeof window === 'undefined') {
      return initialValue
    }

    try {
      // Get from local storage by key
      const item = window.localStorage.getItem(key)
      // Parse stored json or if none return initialValue
      return item ? JSON.parse(item) : initialValue
    } catch (error) {
      // If error also return initialValue
      console.log(error)
      return initialValue
    }
  })

  // Return a wrapped version of useState's setter function that ...
  // ... persists the new value to localStorage.
  const setValue = value => {
    try {
      // Allow value to be a function so we have same API as useState
      const valueToStore = value instanceof Function ? value(storedValue) : value
      // Save state
      setStoredValue(valueToStore)
      // Save to local storage
      window.localStorage.setItem(key, JSON.stringify(valueToStore))
    } catch (error) {
      // A more advanced implementation would handle the error case
      console.log(error)
    }
  }

  return [storedValue, setValue]
}

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
        <br />
        <Divider />
        <br />

        <Typography>
          Minor update
          <ul>
            <li>Added time of recent account banned or suspended on navigation header.</li>
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
