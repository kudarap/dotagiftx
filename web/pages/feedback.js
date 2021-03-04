import React from 'react'
import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import Link from '@material-ui/core/Link'
import Alert from '@material-ui/lab/Alert'
import RedditIcon from '@material-ui/icons/Reddit'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Button from '@/components/Button'
import SteamIcon from '@/components/SteamIcon'
import DiscordIcon from '@/components/DiscordIcon'
import { FormControl, InputLabel, MenuItem, Paper, Select, TextField } from '@material-ui/core'
import {
  REPORT_LABEL_USER_SCAM_ALERT,
  REPORT_TYPE_BUG,
  REPORT_TYPE_FEEDBACK,
  REPORT_TYPE_SCAM_ALERT,
  REPORT_TYPE_SCAM_INCIDENT,
} from '@/constants/report'
import { reportCreate } from '@/service/api'

const useStyles = makeStyles(theme => ({
  main: {
    [theme.breakpoints.down('sm')]: {
      marginTop: theme.spacing(2),
    },
    marginTop: theme.spacing(4),
    // background: 'url("/icon.png") no-repeat bottom right',
    // backgroundSize: 100,
  },
  paper: {
    maxWidth: theme.breakpoints.values.sm,
    margin: '0 auto',
    padding: theme.spacing(2),
  },
}))

export default function About() {
  const classes = useStyles()

  const [payload, setPayload] = React.useState({
    type: REPORT_TYPE_FEEDBACK,
    text: '',
  })

  const [message, setMessage] = React.useState(null)
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(false)

  const handleSubmit = () => {
    setError(null)
    setMessage(null)
    setLoading(true)

    if (payload.text.trim() === '') {
      setError('Description is required')
      setLoading(false)
      return
    }

    ;(async () => {
      try {
        await reportCreate(payload)
        setMessage('Submitted successfully!')
      } catch (e) {
        setError(e.message)
      }

      setLoading(false)
    })()
  }

  const handleSelectChange = e => {
    setPayload({ ...payload, type: e.target.value })
  }

  const handleTextChange = e => {
    setPayload({ ...payload, text: e.target.value })
  }

  return (
    <>
      <Head>
        <title>{APP_NAME} :: Feedback and Report</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Paper className={classes.paper}>
            <Typography variant="h5" component="h1" gutterBottom>
              Feedback and Report
            </Typography>
            {message && (
              <Alert severity="success" variant="filled">
                {message}
              </Alert>
            )}
            <br />
            <form>
              <FormControl fullWidth color="secondary">
                <InputLabel id="demo-simple-select-label">Type</InputLabel>
                <Select
                  labelId="demo-simple-select-label"
                  id="demo-simple-select"
                  value={payload.type}
                  onChange={handleSelectChange}
                  disabled={loading}>
                  <MenuItem value={REPORT_TYPE_FEEDBACK}>Feedback</MenuItem>
                  <MenuItem value={REPORT_TYPE_BUG}>Bug Report</MenuItem>
                  <MenuItem value={REPORT_TYPE_SCAM_ALERT}>Scammer Alert</MenuItem>
                  <MenuItem value={REPORT_TYPE_SCAM_INCIDENT}>Scam Incident</MenuItem>
                </Select>
              </FormControl>
              <br />
              <br />
              <TextField
                fullWidth
                multiline
                label="Description"
                variant="outlined"
                required
                rows={3}
                rowsMax={6}
                color="secondary"
                value={payload.text}
                onInput={handleTextChange}
                helperText="Please paste link from imgur.com if images will be used."
                disabled={loading}
              />
              <br />
              <br />
              <Button
                variant="contained"
                fullWidth
                size="large"
                onClick={handleSubmit}
                disabled={loading}>
                {loading ? 'Submitting...' : 'Submit'}
              </Button>
            </form>
            {error && (
              <Typography align="center" variant="body2" color="error">
                <Typography gutterBottom />
                {error}
              </Typography>
            )}
          </Paper>
        </Container>
      </main>

      <Footer />
    </>
  )
}