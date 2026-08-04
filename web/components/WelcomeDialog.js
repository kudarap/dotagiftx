import React, { useContext, useState, useEffect } from 'react'
import { makeStyles } from 'tss-react/mui'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Typography from '@mui/material/Typography'
import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import Button from '@/components/Button'
import DialogCloseButton from '@/components/DialogCloseButton'
import AppContext from '@/components/AppContext'
import Link from '@/components/Link'
import useLocalStorage from './useLocalStorage'

const storageKey = 'welcome_seen'

const useStyles = makeStyles()(theme => ({
  content: {
    scrollbarWidth: 'thin',
    scrollbarColor: `${theme.palette.divider} transparent`,
    '&::-webkit-scrollbar': {
      width: 8,
    },
    '&::-webkit-scrollbar-track': {
      background: 'transparent',
    },
    '&::-webkit-scrollbar-thumb': {
      backgroundColor: theme.palette.divider,
      borderRadius: 8,
    },
  },
}))

export default function WelcomeDialog() {
  const { classes } = useStyles()
  const { isMobile } = useContext(AppContext)

  const [seen, setSeen] = useLocalStorage(storageKey, false)
  const [open, setOpen] = useState(false)

  useEffect(() => {
    if (seen) {
      return
    }

    const timeoutId = setTimeout(() => {
      setOpen(true)
    }, 1000)

    // Cleanup function to clear the timeout if the component unmounts
    return () => clearTimeout(timeoutId)
  }, [seen])

  const handleClose = () => {
    setOpen(false)
  }

  const handleCheck = evt => {
    setSeen(evt.target.checked)
  }

  return (
    <Dialog
      fullWidth
      fullScreen={isMobile}
      open={open}
      onClose={handleClose}
      aria-labelledby="welcome-dialog-title"
      aria-describedby="welcome-dialog-description">
      <DialogTitle id="welcome-dialog-title">
        Welcome to DotagiftX
        <DialogCloseButton onClick={handleClose} />
      </DialogTitle>
      <DialogContent id="welcome-dialog-description" className={classes.content}>
        <Typography gutterBottom>
          DotagiftX is a peer-to-peer marketplace, trades happen directly between buyer and seller.
          We don&apos;t hold funds or verify identities, so a little caution goes a long way before
          your first trade.
        </Typography>

        <Typography style={{ fontWeight: 'bold' }} gutterBottom>
          If you&apos;re buying
        </Typography>
        <Typography component="div" color="textSecondary" gutterBottom>
          <ul>
            <li>
              Wait for the seller to accept your friend request and reply with a message first.
              Anyone claiming to be the seller without that is a common impersonation scam.
            </li>
            <li>
              Before paying, check the seller&apos;s DotagiftX profile by changing their Steam
              profile URL&apos;s domain to dotagiftx.com, using their Steam ID64 (not a custom URL),
              to see their transaction history and any{' '}
              <Link href="/bans" color="secondary" onClick={handleClose}>
                scam alerts
              </Link>
              .
            </li>
            <li>
              Steam requires 30 days of friendship before an item can be gifted. Double-check the
              seller&apos;s profile again before finalizing.
            </li>
            <li>
              Prefer top sellers or Trader/Partner subscribers. If trading with someone new, paying
              via PayPal Goods &amp; Services offers some dispute protection.
            </li>
            <li>
              For high-value trades, consider using a trusted{' '}
              <Link href="/middleman" color="secondary" onClick={handleClose}>
                middleman
              </Link>
              .
            </li>
            <li>Keep screenshots of your conversation and trade in case of a dispute.</li>
          </ul>
        </Typography>

        <Typography style={{ fontWeight: 'bold' }} gutterBottom>
          If you&apos;re selling
        </Typography>
        <Typography component="div" color="textSecondary" gutterBottom>
          <ul>
            <li>
              Restrict comments on your Steam profile to friends only. This stops scammers from
              posting on your profile while you&apos;re inactive to impersonate you and scam your
              future buyers.
            </li>
          </ul>
        </Typography>

        <Typography>
          Join our{' '}
          <Link
            href="https://discord.gg/UFt9Ny42kM"
            color="secondary"
            target="_blank"
            rel="noreferrer noopener"
            onClick={handleClose}>
            Discord
          </Link>{' '}
          and read the full{' '}
          <Link href="/guides" color="secondary" onClick={handleClose}>
            Guides
          </Link>{' '}
          and{' '}
          <Link href="/rules" color="secondary" onClick={handleClose}>
            Rules
          </Link>{' '}
          before your first trade.
        </Typography>
      </DialogContent>
      <DialogActions>
        <FormControlLabel
          control={<Checkbox onChange={handleCheck} />}
          label={<Typography variant="subtitle2">Don&apos;t show it again</Typography>}
        />
        <Button variant="outlined" color="secondary" onClick={handleClose}>
          Got it
        </Button>
      </DialogActions>
    </Dialog>
  )
}
