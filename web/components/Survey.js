import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import Fab from '@material-ui/core/Fab'
import VoteIcon from '@material-ui/icons/HowToVote'
import VoteDialog from '@/components/VoteDialog'

const useStyles = makeStyles(theme => ({
  root: {},
  fab: {
    position: 'fixed',
    right: theme.spacing(2),
    bottom: theme.spacing(2),
  },
}))

export default function Survey() {
  const classes = useStyles()

  return (
    <div className={classes.root}>
      <VoteDialog open={true} />
      <Fab variant="extended" color="secondary" className={classes.fab}>
        <VoteIcon className={classes.extendedIcon} />
        Vote what&apos;s next
      </Fab>
    </div>
  )
}
