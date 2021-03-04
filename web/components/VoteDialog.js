import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Typography from '@material-ui/core/Typography'
import CircularProgress from '@material-ui/core/CircularProgress'
import Radio from '@material-ui/core/Radio'
import RadioGroup from '@material-ui/core/RadioGroup'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import FormControl from '@material-ui/core/FormControl'
import FormLabel from '@material-ui/core/FormLabel'
import VoteIcon from '@material-ui/icons/HowToVote'
import RemoveIcon from '@material-ui/icons/Close'
import { reportCreate } from '@/service/api'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import { REPORT_LABEL_SURVEY_NEXT, REPORT_TYPE_SURVEY } from '@/constants/report'
import { TextField } from '@material-ui/core'
import { Alert } from '@material-ui/lab'

function shuffle(a) {
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[a[i], a[j]] = [a[j], a[i]]
  }
  return a
}

const voteOptions = shuffle([
  'Inventory import from Steam',
  'Delivery date field on reservation listings',
  'Commend system',
  'Internal messaging',
  'None! everything good',
]).map(v => ({ value: v, label: v }))

export default function VoteDialog(props) {
  const { isMobile } = useContext(AppContext)

  const [notes, setNotes] = React.useState('')
  const [value, setValue] = React.useState('')

  const [message, setMessage] = React.useState('')
  const [error, setError] = React.useState('')
  const [loading, setLoading] = React.useState(false)

  const { onClose } = props
  const handleClose = () => {
    setValue('')
    setError('')
    setLoading(false)
    onClose()
  }

  const handleSubmit = () => {
    const payload = {
      type: REPORT_TYPE_SURVEY,
      label: REPORT_LABEL_SURVEY_NEXT,
      text: value,
    }
    if (value === 'NOA') {
      payload.text = notes
    }

    setLoading(true)
    setError(null)
    ;(async () => {
      try {
        await reportCreate(payload)
        setMessage('Thank you for participating!')
      } catch (e) {
        setError(`Error: ${e.message}`)
      }

      setLoading(false)
    })()
  }

  const handleChange = e => {
    setValue(e.target.value)
  }

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
        Vote what&apos;s next
        <DialogCloseButton onClick={handleClose} />
      </DialogTitle>
      <DialogContent>
        {message && (
          <>
            <Alert
              severity="success"
              variant="filled"
              action={
                <Button size="small" onClick={handleClose}>
                  Close
                </Button>
              }>
              {message}
            </Alert>
            <br />
          </>
        )}
        <Typography>
          Thank you for reaching our latest community goal of <strong>1,000+</strong> items,{' '}
          <strong>500+</strong> users, and <strong>100+</strong> delivered items. Please pick what
          feature we should add next.
        </Typography>
        <br />
        <FormControl component="fieldset">
          <FormLabel component="legend">Here are some suggestions:</FormLabel>
          <RadioGroup aria-label="options" value={value} onChange={handleChange}>
            {voteOptions.map(opts => (
              <FormControlLabel {...opts} control={<Radio />} disabled={message} />
            ))}
            <FormControlLabel
              label="Not listed"
              value="NOA"
              control={<Radio />}
              disabled={message}
            />
          </RadioGroup>
        </FormControl>
        {value === 'NOA' && (
          <TextField
            fullWidth
            label="What is next then?"
            variant="outlined"
            color="secondary"
            value={notes}
            onInput={e => {
              setNotes(e.target.value)
            }}
          />
        )}
      </DialogContent>
      {error && (
        <Typography color="error" align="center" variant="body2">
          {error}
        </Typography>
      )}
      <DialogActions>
        <Button
          disabled={loading}
          startIcon={<RemoveIcon />}
          onClick={handleClose}
          variant="outlined">
          Close
        </Button>
        <Button
          startIcon={loading ? <CircularProgress size={22} color="secondary" /> : <VoteIcon />}
          variant="outlined"
          color="secondary"
          disabled={message}
          onClick={handleSubmit}>
          Submit
        </Button>
      </DialogActions>
    </Dialog>
  )
}
VoteDialog.propTypes = {
  open: PropTypes.bool,
  onClose: PropTypes.func,
}
VoteDialog.defaultProps = {
  open: false,
  onClose: () => {},
}
