import React from 'react'
import { makeStyles } from '@material-ui/core/styles'
import MuiAvatar from '@material-ui/core/Avatar'

const useStyles = makeStyles(() => ({
  root: {
    borderRadius: '15% !important',
  },
}))

export default function Avatar(props) {
  const classes = useStyles()
  return <MuiAvatar classes={{ root: classes.root }} imgProps={{ loading: 'lazy' }} {...props} />
}
