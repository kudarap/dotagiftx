import React, { useContext } from 'react'
import PropTypes from 'prop-types'
import shuffle from 'lodash/shuffle'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Typography from '@mui/material/Typography'
import CircularProgress from '@mui/material/CircularProgress'
import Radio from '@mui/material/Radio'
import RadioGroup from '@mui/material/RadioGroup'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormControl from '@mui/material/FormControl'
import FormLabel from '@mui/material/FormLabel'
import VoteIcon from '@mui/icons-material/HowToVote'
import RemoveIcon from '@mui/icons-material/Close'
import { TextField, Alert } from '@mui/material'
import { reportCreate } from '@/service/api'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import { REPORT_LABEL_SURVEY_NEXT, REPORT_TYPE_SURVEY } from '@/constants/report'

const voteOptions = shuffle([
  'Inventory import from Steam',
  'User commend & report system',
  'Buyer delivery confirmation',
  'Internal messaging',
  'None! everything good',
]).map(v => ({ key: v, value: v, label: v }))

const NONE_OF_THE_ABOVE_OPT = 'NOA'

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
    if (value === NONE_OF_THE_ABOVE_OPT) {
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
              <FormControlLabel {...opts} control={<Radio />} disabled={Boolean(message)} />
            ))}
            <FormControlLabel
              label="Not listed"
              value={NONE_OF_THE_ABOVE_OPT}
              control={<Radio />}
              disabled={Boolean(message)}
            />
          </RadioGroup>
        </FormControl>
        {value === NONE_OF_THE_ABOVE_OPT && (
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
          disabled={Boolean(message)}
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
