import React from 'react'
import PropTypes from 'prop-types'
import { makeStyles } from '@material-ui/core/styles'
import Component from '@material-ui/core/Container'

const useStyles = makeStyles(theme => ({
  root: {
    [theme.breakpoints.down('sm')]: {
      padding: theme.spacing(1),
    },
    minHeight: '40vh',
  },
}))

export default function Container({ children }) {
  const classes = useStyles()

  return (
    <Component className={classes.root} maxWidth="md" disableGutters>
      {children}
    </Component>
  )
}
Container.propTypes = {
  children: PropTypes.node.isRequired,
}
