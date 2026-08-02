import React, { useContext, useEffect } from 'react'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { makeStyles } from 'tss-react/mui'
import Typography from '@mui/material/Typography'
import Alert from '@mui/material/Alert'
import { FormControl, InputLabel, MenuItem, Paper, Select, TextField } from '@mui/material'
import type { SelectChangeEvent } from '@mui/material/Select'
import { APP_NAME } from '@/constants/strings'
import Header from '@/components/Header'
import Container from '@/components/Container'
import Footer from '@/components/Footer'
import Button from '@/components/Button'
import {
  REPORT_TYPE_BUG,
  REPORT_TYPE_FEEDBACK,
  REPORT_TYPE_SCAM_INCIDENT,
} from '@/constants/report'
import { reportCreate } from '@/service/api'
import AppContext from '@/components/AppContext'
import Link from '@/components/Link'

const useStyles = makeStyles()(theme => ({
  main: {
    [theme.breakpoints.down('md')]: {
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

interface Payload {
  type: number
  profile: string
  text: string
  reserved: boolean
}

const defaultPayload: Payload = {
  type: REPORT_TYPE_FEEDBACK,
  profile: '',
  text: '',
  reserved: false,
}

export default function SubmitFeedback() {
  const { classes } = useStyles()

  const [payload, setPayload] = React.useState<Payload>(defaultPayload)

  const [message, setMessage] = React.useState<string | null>(null)
  const [error, setError] = React.useState<string | null>(null)
  const [loading, setLoading] = React.useState(false)

  const { isLoggedIn } = useContext(AppContext)
  const router = useRouter()
  useEffect(() => {
    if (!isLoggedIn) {
      router.push('/')
    }
  }, [isLoggedIn, router])

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
        await reportCreate({
          ...payload,
          text: `${payload.profile} -- ${payload.text}`,
        })
        setMessage('Submitted successfully!')
        setPayload(defaultPayload)
      } catch (e) {
        setError((e as Error).message)
      }

      setLoading(false)
    })()
  }

  const handleSelectChange = (e: SelectChangeEvent<number>) => {
    setPayload({ ...payload, type: e.target.value as number })
  }

  const handleTextChange = (e: React.FormEvent<HTMLDivElement>) => {
    setPayload({ ...payload, text: (e.target as HTMLInputElement).value })
  }

  const handleProfileChange = (e: React.FormEvent<HTMLDivElement>) => {
    setPayload({ ...payload, profile: (e.target as HTMLInputElement).value })
  }

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <title>{APP_NAME} :: Feedback and Report</title>
      </Head>

      <Header />

      <main className={classes.main}>
        <Container>
          <Paper className={classes.paper}>
            <Typography variant="h5" component="h1" gutterBottom>
              Feedback and Report
              <Typography color="textSecondary">
                Feel free to join our{' '}
                <Link
                  color="secondary"
                  target="_blank"
                  rel="noreferrer noopener"
                  href="https://discord.gg/UFt9Ny42kM">
                  Discord
                </Link>{' '}
                if you want to discuss more on your feedback.
              </Typography>
            </Typography>
            {message && (
              <Alert severity="success" variant="filled">
                {message}
              </Alert>
            )}
            <br />
            <form>
              <FormControl fullWidth color="secondary" variant="standard">
                <InputLabel id="demo-simple-select-label">Type</InputLabel>
                <Select
                  labelId="demo-simple-select-label"
                  id="demo-simple-select"
                  value={payload.type}
                  onChange={handleSelectChange}
                  disabled={loading}>
                  <MenuItem value={REPORT_TYPE_FEEDBACK}>Feedback</MenuItem>
                  <MenuItem value={REPORT_TYPE_BUG}>Bug Report</MenuItem>
                  <MenuItem value={REPORT_TYPE_SCAM_INCIDENT}>Scam Incident</MenuItem>
                </Select>
              </FormControl>
              <br />
              <br />
              {payload.type == REPORT_TYPE_SCAM_INCIDENT && (
                <>
                  <TextField
                    sx={{ mb: 1 }}
                    fullWidth
                    required
                    label="Fraudulent Account"
                    variant="outlined"
                    color="secondary"
                    value={payload.profile}
                    onInput={handleProfileChange}
                    helperText="Please paste the link of their Dotagiftx profile NOT Steam profile. We can cross reference the items involved if you have reservation."
                    placeholder="https://dotagiftx.com/profiles/..."
                    disabled={loading}
                  />
                  {/* <FormGroup>
                    <FormControlLabel
                      control={<Checkbox />}
                      label="Is your item reservation?"
                      onChange={handleReservedChange}
                    />
                  </FormGroup> */}
                </>
              )}
              <TextField
                fullWidth
                multiline
                label="Description"
                variant="outlined"
                required
                rows={3}
                maxRows={6}
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
