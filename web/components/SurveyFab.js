import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Fab from '@material-ui/core/Fab'
import VoteIcon from '@material-ui/icons/HowToVote'
import VoteDialog from '@/components/VoteDialog'
import { reportSearch } from '@/service/api'
import { REPORT_TYPE_SURVEY } from '@/constants/report'
import teal from '@material-ui/core/colors/teal'
import { Grow, Slide } from '@material-ui/core'

const useStyles = makeStyles(theme => ({
  root: {},
  fab: {
    position: 'fixed',
    right: theme.spacing(2),
    bottom: theme.spacing(2),
    backgroundColor: teal.A200,
    '&:hover': {
      backgroundColor: teal.A400,
    },
  },
}))

export default function SurveyFab({ userID, label }) {
  const classes = useStyles()

  const [open, setOpen] = React.useState(false)

  // its set true by default so it wont show up in page loads.
  const [voted, setVoted] = React.useState(true)
  React.useEffect(() => {
    if (userID === '') {
      return
    }

    ;(async () => {
      const res = await reportSearch({
        user_id: userID,
        type: REPORT_TYPE_SURVEY,
        label,
      })
      if (res && res.result_count === 0) {
        setTimeout(() => {
          setVoted(false)
        }, 2310) // average time on page base on GA
      }
    })()
  }, [userID])

  const handleClose = () => {
    setOpen(false)
    setVoted(true)
  }

  // if (voted) {
  //   return null
  // }

  return (
    <div className={classes.root}>
      <VoteDialog open={open} onClose={handleClose} />
      <Slide direction="up" in={!voted} mountOnEnter unmountOnExit>
        <Fab variant="extended" className={classes.fab} onClick={() => setOpen(true)}>
          <VoteIcon className={classes.extendedIcon} />
          Vote what&apos;s next
        </Fab>
      </Slide>
    </div>
  )
}
SurveyFab.propTypes = {
  userID: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
}
